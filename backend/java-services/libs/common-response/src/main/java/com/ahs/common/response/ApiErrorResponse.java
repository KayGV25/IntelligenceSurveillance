package com.ahs.common.response;

public record ApiErrorResponse(
        String code,
        String key,
        String message
) {
}
