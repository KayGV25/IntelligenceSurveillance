package com.ahs.common.observability;

public final class TraceConstants {

    private TraceConstants() {
    }

    public static final String TRACE_ID_HEADER = "X-Trace-Id";
    public static final String REQUEST_ID_HEADER = "X-Request-Id";
    public static final String CORRELATION_ID_HEADER = "X-Correlation-Id";

    public static final String TRACE_ID_ATTRIBUTE = "traceId";
    public static final String REQUEST_ID_ATTRIBUTE = "requestId";
    public static final String CORRELATION_ID_ATTRIBUTE = "correlationId";
}