package com.ahs.api_gateway.exception;

import com.ahs.common.error.ErrorCode;
import com.ahs.common.error.ErrorKey;
import com.ahs.common.response.ApiResponse;
import lombok.RequiredArgsConstructor;
import org.springframework.boot.webflux.error.ErrorWebExceptionHandler;
import org.springframework.core.annotation.Order;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Component;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.web.server.ServerWebExchange;
import reactor.core.publisher.Mono;
import tools.jackson.databind.ObjectMapper;

@Component
@Order(-2)
@RequiredArgsConstructor
public class GlobalExceptionHandler implements ErrorWebExceptionHandler {

    private final ErrorResponseWriter errorResponseWriter;

    @Override
    public Mono<Void> handle(
            ServerWebExchange exchange,
            Throwable ex
    ) {
        HttpStatus status = HttpStatus.INTERNAL_SERVER_ERROR;
        ErrorCode code = ErrorCode.INTERNAL_ERROR;
        ErrorKey key = ErrorKey.INTERNAL_SERVER_ERROR;
        String message = "Internal server error";

        if (ex instanceof ResponseStatusException responseStatusException) {
            status = HttpStatus.valueOf(
                    responseStatusException.getStatusCode().value()
            );

            if (status == HttpStatus.NOT_FOUND) {
                code = ErrorCode.NOT_FOUND;
                key = ErrorKey.ROUTE_NOT_FOUND;
                message = "Route not found";
            } else if (status == HttpStatus.UNAUTHORIZED) {
                code = ErrorCode.UNAUTHORIZED;
                key = ErrorKey.INSUFFICIENT_PERMISSION;
                message = "Unauthorized";
            } else if (status == HttpStatus.FORBIDDEN) {
                code = ErrorCode.FORBIDDEN;
                key = ErrorKey.INSUFFICIENT_PERMISSION;
                message = "Forbidden";
            } else if (status == HttpStatus.BAD_REQUEST) {
                code = ErrorCode.BAD_REQUEST;
                key = ErrorKey.INVALID_REQUEST;
                message = "Bad request";
            } else if (status == HttpStatus.SERVICE_UNAVAILABLE) {
                code = ErrorCode.SERVICE_UNAVAILABLE;
                key = ErrorKey.SERVICE_UNAVAILABLE;
                message = "Service unavailable";
            } else if (status == HttpStatus.TOO_MANY_REQUESTS) {
                code = ErrorCode.RATE_LIMITED;
                key = ErrorKey.RATE_LIMIT_EXCEEDED;
                message = "Too many requests";
            }
        }

        return errorResponseWriter.write(exchange, status, code, key, message);
    }
}
