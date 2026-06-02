package com.ahs.common.error;

public class BusinessException extends RuntimeException {

    private final ErrorCode errorCode;
    private final ErrorKey errorKey;

    public BusinessException(
            ErrorCode errorCode,
            ErrorKey errorKey,
            String message
    ) {
        super(message);
        this.errorCode = errorCode;
        this.errorKey = errorKey;
    }

    public ErrorCode getErrorCode() {
        return errorCode;
    }

    public ErrorKey getErrorKey() {
        return errorKey;
    }
}
