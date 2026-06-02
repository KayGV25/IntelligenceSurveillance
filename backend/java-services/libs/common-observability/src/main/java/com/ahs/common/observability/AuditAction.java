package com.ahs.common.observability;

public enum AuditAction {
    GATEWAY_REQUEST,
    GATEWAY_REQUEST_FAILED,
    GATEWAY_RATE_LIMITED,
    AUTHENTICATION_FAILED,
    AUTHORIZATION_FAILED
}
