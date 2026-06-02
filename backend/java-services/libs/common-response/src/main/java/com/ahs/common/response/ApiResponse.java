package com.ahs.common.response;

import java.time.Instant;

public record ApiResponse<T>(
        boolean success,
        T data,
        ApiErrorResponse error,
        Instant timestamp
) {
        public static <T> ApiResponse<T> success(T data) {
            return new ApiResponse<T>(
                    true,
                    data,
                    null,
                    Instant.now()
            );
        }

        public static <T> ApiResponse<T> failure(
                String code,
                String key,
                String message
        ) {
            return new ApiResponse<T>(
                    false,
                    null,
                    new ApiErrorResponse(code, key, message),
                    Instant.now()
            );
        }
}