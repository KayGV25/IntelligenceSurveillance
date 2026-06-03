package com.ahs.auth_service.exception;

import com.ahs.common.error.BusinessException;
import com.ahs.common.error.ErrorCode;
import com.ahs.common.error.ErrorKey;
import com.ahs.common.response.ApiResponse;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;
import reactor.core.publisher.Mono;

@RestControllerAdvice
public class AuthExceptionHandler {

    @ExceptionHandler(BusinessException.class)
    public Mono<ResponseEntity<ApiResponse<Void>>> handleBusinessException(BusinessException ex) {
        return Mono.just(ResponseEntity
                .status(resolveStatus(ex.getErrorCode()))
                .body(ApiResponse.failure(
                        ex.getErrorCode().name(),
                        ex.getErrorKey().name(),
                        ex.getMessage()
                )));
    }

    @ExceptionHandler(Exception.class)
    public Mono<ResponseEntity<ApiResponse<Void>>> handleException(Exception ex) {
        return Mono.just(ResponseEntity
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(ApiResponse.failure(
                        ErrorCode.INTERNAL_ERROR.name(),
                        ErrorKey.INTERNAL_SERVER_ERROR.name(),
                        "Internal server error"
                )));
    }

    private HttpStatus resolveStatus(ErrorCode code) {
        return switch (code) {
            case BAD_REQUEST, VALIDATION_ERROR -> HttpStatus.BAD_REQUEST;
            case UNAUTHORIZED -> HttpStatus.UNAUTHORIZED;
            case FORBIDDEN -> HttpStatus.FORBIDDEN;
            case NOT_FOUND -> HttpStatus.NOT_FOUND;
            case CONFLICT -> HttpStatus.CONFLICT;
            case RATE_LIMITED -> HttpStatus.TOO_MANY_REQUESTS;
            case SERVICE_UNAVAILABLE -> HttpStatus.SERVICE_UNAVAILABLE;
            default -> HttpStatus.INTERNAL_SERVER_ERROR;
        };
    }
}
