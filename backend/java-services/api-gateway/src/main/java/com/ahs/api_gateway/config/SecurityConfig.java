package com.ahs.api_gateway.config;

import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.config.web.server.ServerHttpSecurity;
import org.springframework.security.web.server.SecurityWebFilterChain;

@Configuration
@RequiredArgsConstructor
public class SecurityConfig {

    private final SecurityProperties securityProperties;

    @Bean
    public SecurityWebFilterChain securityWebFilterChain(ServerHttpSecurity http) {
        http
                .csrf(ServerHttpSecurity.CsrfSpec::disable)
                .authorizeExchange(exchange -> {
                    exchange
                            .pathMatchers(HttpMethod.OPTIONS, "/**")
                            .permitAll();
                    exchange
                            .pathMatchers(
                                    "/actuator/health",
                                    "/actuator/info",
                                    "/api/v1/gateway/health",
                                    "/api/v1/auth/**",
                                    "/v3/api-docs/**",
                                    "/swagger-ui.html",
                                    "/swagger-ui/**",
                                    "/webjars/**"
                            )
                            .permitAll();

                    if(securityProperties.jwtEnabled()) {
                        exchange.anyExchange().authenticated();
                    } else {
                        exchange.anyExchange().permitAll();
                    }
                });

        if (securityProperties.jwtEnabled()) {
            http.oauth2ResourceServer(oauth2 -> oauth2.jwt(jwt -> {

            }));
        }
        return http.build();
    }
}
