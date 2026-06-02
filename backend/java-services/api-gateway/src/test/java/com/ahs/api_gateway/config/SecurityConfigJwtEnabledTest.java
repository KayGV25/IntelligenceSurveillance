package com.ahs.api_gateway.config;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.context.annotation.Bean;
import org.springframework.security.oauth2.jwt.ReactiveJwtDecoder;
import org.springframework.security.web.server.SecurityWebFilterChain;
import reactor.core.publisher.Mono;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest(properties = {
        "ahs.security.jwt-enabled=true",
        "ahs.services.auth-service-url=http://auth-service:8081",
        "ahs.services.camera-service-url=http://camera-service:8101",
        "ahs.services.map-service-url=http://map-service:8105",
        "ahs.services.alert-service-url=http://alert-service:8106",
        "ahs.services.incident-service-url=http://incident-service:8107",
        "ahs.services.tracking-service-url=http://tracking-service:8108",
        "ahs.services.playback-service-url=http://playback-service:8104"
})
class SecurityConfigJwtEnabledTest {

    @Autowired
    private SecurityWebFilterChain securityWebFilterChain;

    @Test
    void buildsResourceServerSecurityFilterChainWhenJwtIsEnabled() {
        assertThat(securityWebFilterChain).isNotNull();
    }

    @TestConfiguration
    static class JwtDecoderConfig {

        @Bean
        ReactiveJwtDecoder reactiveJwtDecoder() {
            return token -> Mono.error(new UnsupportedOperationException("Test decoder"));
        }
    }
}
