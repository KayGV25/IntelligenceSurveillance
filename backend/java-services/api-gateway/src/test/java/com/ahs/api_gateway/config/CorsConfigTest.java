package com.ahs.api_gateway.config;

import org.junit.jupiter.api.Test;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.mock.http.server.reactive.MockServerHttpRequest;
import org.springframework.mock.web.server.MockServerWebExchange;
import org.springframework.web.cors.reactive.CorsWebFilter;
import org.springframework.web.server.WebFilterChain;
import reactor.core.publisher.Mono;

import java.util.concurrent.atomic.AtomicBoolean;

import static org.assertj.core.api.Assertions.assertThat;

class CorsConfigTest {

    @Test
    void corsFilterAllowsLocalhostPreflightAndExposesTraceHeaders() {
        CorsWebFilter filter = new CorsConfig().corsFilter();
        MockServerWebExchange exchange = MockServerWebExchange.from(
                MockServerHttpRequest.options("http://api-gateway.local/api/v1/gateway/health")
                        .header(HttpHeaders.ORIGIN, "http://localhost:5173")
                        .header(HttpHeaders.ACCESS_CONTROL_REQUEST_METHOD, HttpMethod.GET.name())
        );
        AtomicBoolean chainCalled = new AtomicBoolean(false);
        WebFilterChain chain = ignored -> {
            chainCalled.set(true);
            return Mono.empty();
        };

        filter.filter(exchange, chain).block();

        assertThat(chainCalled).isFalse();
        assertThat(exchange.getResponse().getHeaders().getAccessControlAllowOrigin())
                .isEqualTo("http://localhost:5173");
        assertThat(exchange.getResponse().getHeaders().getAccessControlAllowMethods())
                .contains(HttpMethod.GET, HttpMethod.POST, HttpMethod.PUT, HttpMethod.PATCH, HttpMethod.DELETE, HttpMethod.OPTIONS);
        assertThat(exchange.getResponse().getHeaders().getAccessControlExposeHeaders())
                .contains("X-Trace-Id", "X-Request-Id", "X-Correlation-Id");
        assertThat(exchange.getResponse().getHeaders().getAccessControlMaxAge()).isEqualTo(3600L);
    }
}
