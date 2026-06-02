package com.ahs.api_gateway.config;

import org.junit.jupiter.api.Test;
import org.springframework.cloud.gateway.filter.ratelimit.KeyResolver;
import org.springframework.cloud.gateway.filter.ratelimit.RedisRateLimiter;
import org.springframework.mock.http.server.reactive.MockServerHttpRequest;
import org.springframework.mock.web.server.MockServerWebExchange;
import reactor.test.StepVerifier;

import java.net.InetSocketAddress;

import static org.assertj.core.api.Assertions.assertThat;

class RateLimitConfigTest {

    private final RateLimitConfig config = new RateLimitConfig();

    @Test
    void createsRedisRateLimiterBean() {
        RedisRateLimiter limiter = config.redisRateLimiter();

        assertThat(limiter).isNotNull();
    }

    @Test
    void ipKeyResolverUsesRemoteIpAddress() {
        KeyResolver resolver = config.ipKeyResolver();
        MockServerWebExchange exchange = MockServerWebExchange.from(
                MockServerHttpRequest.get("/api/v1/cameras")
                        .remoteAddress(new InetSocketAddress("192.168.1.50", 4321))
        );

        StepVerifier.create(resolver.resolve(exchange))
                .expectNext("192.168.1.50")
                .verifyComplete();
    }

    @Test
    void ipKeyResolverFallsBackToUnknownWhenRemoteAddressMissing() {
        KeyResolver resolver = config.ipKeyResolver();
        MockServerWebExchange exchange = MockServerWebExchange.from(
                MockServerHttpRequest.get("/api/v1/cameras")
        );

        StepVerifier.create(resolver.resolve(exchange))
                .expectNext("unknown")
                .verifyComplete();
    }
}
