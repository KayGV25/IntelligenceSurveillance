package com.ahs.api_gateway.exception;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.mock.http.server.reactive.MockServerHttpRequest;
import org.springframework.mock.web.server.MockServerWebExchange;
import org.springframework.web.server.ResponseStatusException;
import tools.jackson.databind.ObjectMapper;

import java.util.stream.Stream;

import static org.assertj.core.api.Assertions.assertThat;

class GlobalExceptionHandlerTest {

    private final GlobalExceptionHandler handler = new GlobalExceptionHandler(
            new ErrorResponseWriter(new ObjectMapper())
    );

    @ParameterizedTest
    @MethodSource("gatewayErrors")
    void responseStatusExceptionsAreMappedToApiErrorResponses(
            HttpStatus status,
            String code,
            String key,
            String message
    ) {
        MockServerWebExchange exchange = MockServerWebExchange.from(
                MockServerHttpRequest.get("/missing")
        );

        handler.handle(exchange, new ResponseStatusException(status)).block();

        assertThat(exchange.getResponse().getStatusCode()).isEqualTo(status);
        assertThat(exchange.getResponse().getHeaders().getContentType()).isEqualTo(MediaType.APPLICATION_JSON);
        assertThat(exchange.getResponse().getBodyAsString().block())
                .contains("\"success\":false")
                .contains("\"error\":{\"code\":\"" + code + "\",\"key\":\"" + key + "\",\"message\":\"" + message + "\"")
                .contains("\"timestamp\":");
    }

    @Test
    void unmappedExceptionsReturnInternalErrorResponse() {
        MockServerWebExchange exchange = MockServerWebExchange.from(
                MockServerHttpRequest.get("/boom")
        );

        handler.handle(exchange, new IllegalStateException("boom")).block();

        assertThat(exchange.getResponse().getStatusCode()).isEqualTo(HttpStatus.INTERNAL_SERVER_ERROR);
        assertThat(exchange.getResponse().getHeaders().getContentType()).isEqualTo(MediaType.APPLICATION_JSON);
        assertThat(exchange.getResponse().getBodyAsString().block())
                .contains("\"success\":false")
                .contains("\"error\":{\"code\":\"INTERNAL_ERROR\",\"key\":\"INTERNAL_SERVER_ERROR\",\"message\":\"Internal server error\"")
                .contains("\"timestamp\":");
    }

    private static Stream<Arguments> gatewayErrors() {
        return Stream.of(
                Arguments.of(HttpStatus.NOT_FOUND, "NOT_FOUND", "ROUTE_NOT_FOUND", "Route not found"),
                Arguments.of(HttpStatus.UNAUTHORIZED, "UNAUTHORIZED", "INSUFFICIENT_PERMISSION", "Unauthorized"),
                Arguments.of(HttpStatus.FORBIDDEN, "FORBIDDEN", "INSUFFICIENT_PERMISSION", "Forbidden"),
                Arguments.of(HttpStatus.BAD_REQUEST, "BAD_REQUEST", "INVALID_REQUEST", "Bad request"),
                Arguments.of(HttpStatus.SERVICE_UNAVAILABLE, "SERVICE_UNAVAILABLE", "SERVICE_UNAVAILABLE", "Service unavailable"),
                Arguments.of(HttpStatus.TOO_MANY_REQUESTS, "RATE_LIMITED", "RATE_LIMIT_EXCEEDED", "Too many requests")
        );
    }
}
