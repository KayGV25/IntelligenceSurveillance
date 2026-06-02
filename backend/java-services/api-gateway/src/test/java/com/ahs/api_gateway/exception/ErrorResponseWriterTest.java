package com.ahs.api_gateway.exception;

import com.ahs.common.error.ErrorCode;
import com.ahs.common.error.ErrorKey;
import org.junit.jupiter.api.Test;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.mock.http.server.reactive.MockServerHttpRequest;
import org.springframework.mock.web.server.MockServerWebExchange;
import tools.jackson.databind.ObjectMapper;

import static org.assertj.core.api.Assertions.assertThat;

class ErrorResponseWriterTest {

    private final ErrorResponseWriter writer = new ErrorResponseWriter(new ObjectMapper());

    @Test
    void writesSharedApiErrorResponseBody() {
        MockServerWebExchange exchange = MockServerWebExchange.from(
                MockServerHttpRequest.get("/secure")
        );

        writer.write(
                exchange,
                HttpStatus.UNAUTHORIZED,
                ErrorCode.UNAUTHORIZED,
                ErrorKey.AUTHENTICATION_FAILED,
                "Authentication required"
        ).block();

        assertThat(exchange.getResponse().getStatusCode()).isEqualTo(HttpStatus.UNAUTHORIZED);
        assertThat(exchange.getResponse().getHeaders().getContentType()).isEqualTo(MediaType.APPLICATION_JSON);
        assertThat(exchange.getResponse().getBodyAsString().block())
                .contains("\"success\":false")
                .contains("\"error\":{\"code\":\"UNAUTHORIZED\",\"key\":\"AUTHENTICATION_FAILED\",\"message\":\"Authentication required\"")
                .contains("\"timestamp\":");
    }
}
