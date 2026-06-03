package com.ahs.auth_service.audit;

import com.ahs.common.observability.RequestIdGenerator;
import com.ahs.common.observability.TraceConstants;
import org.springframework.http.server.reactive.ServerHttpRequest;

public record RequestAuditContext(
        String ipAddress,
        String userAgent,
        String traceId,
        String requestId,
        String correlationId
) {
    public static RequestAuditContext from(ServerHttpRequest request) {
        return new RequestAuditContext(
                request.getRemoteAddress() == null ? "unknown" : request.getRemoteAddress().toString(),
                request.getHeaders().getFirst("User-Agent"),
                getOrCreateHeader(request, TraceConstants.TRACE_ID_HEADER),
                getOrCreateHeader(request, TraceConstants.REQUEST_ID_HEADER),
                getOrCreateHeader(request, TraceConstants.CORRELATION_ID_HEADER)
        );
    }

    private static String getOrCreateHeader(
            ServerHttpRequest request,
            String headerName
    ) {
        String value = request.getHeaders().getFirst(headerName);
        return value == null || value.isBlank()
                ? RequestIdGenerator.generate()
                : value;
    }
}
