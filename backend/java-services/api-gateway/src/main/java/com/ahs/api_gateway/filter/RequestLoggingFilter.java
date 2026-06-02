package com.ahs.api_gateway.filter;

import com.ahs.api_gateway.context.GatewayRequestContext;
import com.ahs.common.observability.TraceConstants;
import lombok.extern.slf4j.Slf4j;
import org.slf4j.MDC;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.http.server.reactive.ServerHttpRequest;
import org.springframework.stereotype.Component;
import org.springframework.web.server.ServerWebExchange;
import org.springframework.web.server.WebFilter;
import org.springframework.web.server.WebFilterChain;
import reactor.core.publisher.Mono;

import java.time.Instant;
import java.util.UUID;

@Slf4j
@Component
@Order(-100)
public class RequestLoggingFilter implements WebFilter, Ordered {

    @Override
    public Mono<Void> filter(
            ServerWebExchange exchange,
            WebFilterChain chain
    ) {
        long startTimeMs = System.currentTimeMillis();

        ServerHttpRequest request = exchange.getRequest();

        String traceId = getOrCreateHeader(request, TraceConstants.TRACE_ID_HEADER);
        String requestId = getOrCreateHeader(request, TraceConstants.REQUEST_ID_HEADER);
        String correlationId = getOrCreateHeader(request, TraceConstants.CORRELATION_ID_HEADER);

        GatewayRequestContext requestContext = new GatewayRequestContext(
                traceId,
                requestId,
                correlationId,
                request.getMethod().name(),
                request.getURI().getPath(),
                request.getRemoteAddress() == null ? "unknown" : request.getRemoteAddress().toString(),
                Instant.now()
        );

        ServerHttpRequest mutatedRequest = request.mutate()
                .header(TraceConstants.TRACE_ID_HEADER, traceId)
                .header(TraceConstants.REQUEST_ID_HEADER, requestId)
                .header(TraceConstants.CORRELATION_ID_HEADER, correlationId)
                .build();

        exchange.getAttributes().put("gatewayRequestContext", requestContext);
        exchange.getAttributes().put(TraceConstants.TRACE_ID_ATTRIBUTE, traceId);
        exchange.getAttributes().put(TraceConstants.REQUEST_ID_ATTRIBUTE, requestId);
        exchange.getAttributes().put(TraceConstants.CORRELATION_ID_ATTRIBUTE, correlationId);

        exchange.getResponse().beforeCommit(() -> {
            exchange.getResponse().getHeaders().set(TraceConstants.TRACE_ID_HEADER, traceId);
            exchange.getResponse().getHeaders().set(TraceConstants.REQUEST_ID_HEADER, requestId);
            exchange.getResponse().getHeaders().set(TraceConstants.CORRELATION_ID_HEADER, correlationId);
            return Mono.empty();
        });

        return Mono.defer(() -> {
                    putMdc(traceId, requestId, correlationId);

                    log.info(
                            "gateway_request_started method={} path={} traceId={} requestId={} correlationId={} remoteAddress={}",
                            requestContext.method(),
                            requestContext.path(),
                            traceId,
                            requestId,
                            correlationId,
                            requestContext.remoteAddress()
                    );

                    return chain.filter(exchange.mutate().request(mutatedRequest).build());
                })
                .doFinally(signalType -> {
                    putMdc(traceId, requestId, correlationId);

                    long durationMs = System.currentTimeMillis() - startTimeMs;

                    log.info(
                            "gateway_request_completed method={} path={} status={} durationMs={} traceId={} requestId={} correlationId={}",
                            requestContext.method(),
                            requestContext.path(),
                            exchange.getResponse().getStatusCode(),
                            durationMs,
                            traceId,
                            requestId,
                            correlationId
                    );

                    MDC.clear();
                });
    }

    private String getOrCreateHeader(ServerHttpRequest request, String headerName) {
        String value = request.getHeaders().getFirst(headerName);
        return value == null || value.isBlank()
                ? UUID.randomUUID().toString()
                : value;
    }

    private void putMdc(
            String traceId,
            String requestId,
            String correlationId
    ) {
        MDC.put("traceId", traceId);
        MDC.put("requestId", requestId);
        MDC.put("correlationId", correlationId);
    }

    @Override
    public int getOrder() {
        return -100;
    }
}
