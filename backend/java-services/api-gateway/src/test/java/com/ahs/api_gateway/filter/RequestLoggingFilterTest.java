package com.ahs.api_gateway.filter;

import com.ahs.api_gateway.context.GatewayRequestContext;
import com.ahs.common.observability.TraceConstants;
import org.junit.jupiter.api.Test;
import org.slf4j.MDC;
import org.springframework.mock.http.server.reactive.MockServerHttpRequest;
import org.springframework.mock.web.server.MockServerWebExchange;
import org.springframework.web.server.WebFilterChain;
import reactor.core.publisher.Mono;

import java.net.InetSocketAddress;
import java.util.UUID;
import java.util.concurrent.atomic.AtomicReference;

import static org.assertj.core.api.Assertions.assertThat;

class RequestLoggingFilterTest {

    private final RequestLoggingFilter filter = new RequestLoggingFilter();

    @Test
    void preservesIncomingCorrelationHeadersAndAddsContext() {
        MockServerWebExchange exchange = MockServerWebExchange.from(
                MockServerHttpRequest.get("/api/v1/cameras")
                        .remoteAddress(new InetSocketAddress("10.0.0.5", 12345))
                        .header(TraceConstants.TRACE_ID_HEADER, "trace-1")
                        .header(TraceConstants.REQUEST_ID_HEADER, "request-1")
                        .header(TraceConstants.CORRELATION_ID_HEADER, "correlation-1")
        );
        AtomicReference<String> forwardedTraceId = new AtomicReference<>();
        WebFilterChain chain = filteredExchange -> {
            forwardedTraceId.set(filteredExchange.getRequest().getHeaders().getFirst(TraceConstants.TRACE_ID_HEADER));
            return filteredExchange.getResponse().setComplete();
        };

        filter.filter(exchange, chain).block();

        assertThat(forwardedTraceId).hasValue("trace-1");
        assertThat(exchange.getAttributes())
                .containsEntry(TraceConstants.TRACE_ID_ATTRIBUTE, "trace-1")
                .containsEntry(TraceConstants.REQUEST_ID_ATTRIBUTE, "request-1")
                .containsEntry(TraceConstants.CORRELATION_ID_ATTRIBUTE, "correlation-1");
        assertThat(exchange.getResponse().getHeaders().getFirst(TraceConstants.TRACE_ID_HEADER)).isEqualTo("trace-1");
        assertThat(exchange.getResponse().getHeaders().getFirst(TraceConstants.REQUEST_ID_HEADER)).isEqualTo("request-1");
        assertThat(exchange.getResponse().getHeaders().getFirst(TraceConstants.CORRELATION_ID_HEADER)).isEqualTo("correlation-1");
        Object gatewayRequestContext = exchange.getAttribute("gatewayRequestContext");

        assertThat(gatewayRequestContext)
                .isInstanceOfSatisfying(GatewayRequestContext.class, context -> {
                    assertThat(context.traceId()).isEqualTo("trace-1");
                    assertThat(context.requestId()).isEqualTo("request-1");
                    assertThat(context.correlationId()).isEqualTo("correlation-1");
                    assertThat(context.method()).isEqualTo("GET");
                    assertThat(context.path()).isEqualTo("/api/v1/cameras");
                    assertThat(context.remoteAddress()).contains("10.0.0.5");
                    assertThat(context.startTime()).isNotNull();
                });
        assertThat(MDC.get("traceId")).isNull();
    }

    @Test
    void generatesCorrelationHeadersWhenMissing() {
        MockServerWebExchange exchange = MockServerWebExchange.from(
                MockServerHttpRequest.post("/api/v1/alerts")
        );
        AtomicReference<String> forwardedTraceId = new AtomicReference<>();
        AtomicReference<String> forwardedRequestId = new AtomicReference<>();
        AtomicReference<String> forwardedCorrelationId = new AtomicReference<>();
        WebFilterChain chain = filteredExchange -> {
            forwardedTraceId.set(filteredExchange.getRequest().getHeaders().getFirst(TraceConstants.TRACE_ID_HEADER));
            forwardedRequestId.set(filteredExchange.getRequest().getHeaders().getFirst(TraceConstants.REQUEST_ID_HEADER));
            forwardedCorrelationId.set(filteredExchange.getRequest().getHeaders().getFirst(TraceConstants.CORRELATION_ID_HEADER));
            return filteredExchange.getResponse().setComplete();
        };

        filter.filter(exchange, chain).block();

        assertThat(UUID.fromString(forwardedTraceId.get())).isNotNull();
        assertThat(UUID.fromString(forwardedRequestId.get())).isNotNull();
        assertThat(UUID.fromString(forwardedCorrelationId.get())).isNotNull();
        assertThat(exchange.getResponse().getHeaders().getFirst(TraceConstants.TRACE_ID_HEADER)).isEqualTo(forwardedTraceId.get());
        assertThat(exchange.getResponse().getHeaders().getFirst(TraceConstants.REQUEST_ID_HEADER)).isEqualTo(forwardedRequestId.get());
        assertThat(exchange.getResponse().getHeaders().getFirst(TraceConstants.CORRELATION_ID_HEADER)).isEqualTo(forwardedCorrelationId.get());
    }

    @Test
    void hasHighestGatewayFilterOrder() {
        assertThat(filter.getOrder()).isEqualTo(-100);
    }
}
