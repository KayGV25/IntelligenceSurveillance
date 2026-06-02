package com.ahs.common.observability;

import java.time.Instant;
import java.util.Map;

public record AuditEvent(
        AuditAction action,
        String service,
        String traceId,
        String correlationId,
        String userId,
        String method,
        String path,
        String ipAddress,
        Instant timestamp,
        Map<String, Object> metadata
) {
}
