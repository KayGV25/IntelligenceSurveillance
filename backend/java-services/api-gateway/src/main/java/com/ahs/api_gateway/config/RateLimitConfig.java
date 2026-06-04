package com.ahs.api_gateway.config;

import org.springframework.cloud.gateway.filter.ratelimit.KeyResolver;
import org.springframework.cloud.gateway.filter.ratelimit.RedisRateLimiter;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Primary;
import reactor.core.publisher.Mono;

@Configuration
public class RateLimitConfig {

    @Bean
    @Primary
    public RedisRateLimiter redisRateLimiter() {
        RedisRateLimiter limiter = new RedisRateLimiter(10, 20);
        limiter.setIncludeHeaders(false);
        return limiter;
    }

    @Bean
    public RedisRateLimiter loginRedisRateLimiter() {
        RedisRateLimiter limiter = new RedisRateLimiter(1, 5);
        limiter.setIncludeHeaders(false);
        return limiter;
    }

    @Bean
    public KeyResolver ipKeyResolver() {
        return exchange -> {
            if (exchange.getRequest().getRemoteAddress() == null) {
                return Mono.just("unknown");
            }

            return Mono.just(
                    exchange.getRequest()
                            .getRemoteAddress()
                            .getAddress().getHostAddress()
            );
        };
    }
}
