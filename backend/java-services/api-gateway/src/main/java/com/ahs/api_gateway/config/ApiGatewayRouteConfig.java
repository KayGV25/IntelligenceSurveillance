package com.ahs.api_gateway.config;

import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.cloud.gateway.filter.ratelimit.KeyResolver;
import org.springframework.cloud.gateway.filter.ratelimit.RedisRateLimiter;
import org.springframework.cloud.gateway.route.RouteLocator;
import org.springframework.cloud.gateway.route.builder.RouteLocatorBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class ApiGatewayRouteConfig {

    private final ApiGatewayProperties properties;
    private final RedisRateLimiter redisRateLimiter;
    private final RedisRateLimiter loginRedisRateLimiter;
    private final KeyResolver resolver;

    public ApiGatewayRouteConfig(
            ApiGatewayProperties properties,
            @Qualifier("redisRateLimiter") RedisRateLimiter redisRateLimiter,
            @Qualifier("loginRedisRateLimiter") RedisRateLimiter loginRedisRateLimiter,
            KeyResolver resolver
    ) {
        this.properties = properties;
        this.redisRateLimiter = redisRateLimiter;
        this.loginRedisRateLimiter = loginRedisRateLimiter;
        this.resolver = resolver;
    }

    @Bean
    public RouteLocator apiGatewayRoutes(RouteLocatorBuilder builder) {
        return builder.routes()
                .route("auth-service", route -> route
                        .path("/api/v1/auth/**")
                        .filters(filter -> filter.requestRateLimiter(config -> config
                                .setRateLimiter(loginRedisRateLimiter)
                                .setKeyResolver(resolver)))
                        .uri(properties.authServiceUrl()))

                .route("camera-service", route -> route
                        .path("/api/v1/cameras/**")
                        .filters(filter -> filter.requestRateLimiter(config -> config
                                .setRateLimiter(redisRateLimiter)
                                .setKeyResolver(resolver)))
                        .uri(properties.cameraServiceUrl()))

                .route("map-service", route -> route
                        .path("/api/v1/maps/**")
                        .filters(filter -> filter.requestRateLimiter(config -> config
                                .setRateLimiter(redisRateLimiter)
                                .setKeyResolver(resolver)))
                        .uri(properties.mapServiceUrl()))

                .route("alert-service", route -> route
                        .path("/api/v1/alerts/**")
                        .filters(filter -> filter.requestRateLimiter(config -> config
                                .setRateLimiter(redisRateLimiter)
                                .setKeyResolver(resolver)))
                        .uri(properties.alertServiceUrl()))

                .route("incident-service", route -> route
                        .path("/api/v1/incidents/**")
                        .filters(filter -> filter.requestRateLimiter(config -> config
                                .setRateLimiter(redisRateLimiter)
                                .setKeyResolver(resolver)))
                        .uri(properties.incidentServiceUrl()))

                .route("tracking-service", route -> route
                        .path("/api/v1/tracking/**")
                        .filters(filter -> filter.requestRateLimiter(config -> config
                                .setRateLimiter(redisRateLimiter)
                                .setKeyResolver(resolver)))
                        .uri(properties.trackingServiceUrl()))

                .route("playback-service", route -> route
                        .path("/api/v1/playback/**")
                        .filters(filter -> filter.requestRateLimiter(config -> config
                                .setRateLimiter(redisRateLimiter)
                                .setKeyResolver(resolver)))
                        .uri(properties.playbackServiceUrl()))
                .build();
    }
}
