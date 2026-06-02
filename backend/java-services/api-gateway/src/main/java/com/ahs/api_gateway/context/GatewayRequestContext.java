package com.ahs.api_gateway.context;

import java.time.Instant;

public record GatewayRequestContext(
        String traceId,
        String requestId,
        String correlationId,
        String method,
        String path,
        String remoteAddress,
        Instant startTime
) {
}
