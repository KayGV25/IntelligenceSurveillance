package com.ahs.api_gateway.config;

import com.ahs.api_gateway.exception.ErrorResponseWriter;
import com.ahs.common.error.ErrorCode;
import com.ahs.common.error.ErrorKey;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.security.config.web.server.ServerHttpSecurity;
import org.springframework.security.web.server.SecurityWebFilterChain;

@Configuration
@RequiredArgsConstructor
public class SecurityConfig {

    private final SecurityProperties securityProperties;
    private final ErrorResponseWriter errorResponseWriter;

    @Bean
    public SecurityWebFilterChain securityWebFilterChain(ServerHttpSecurity http) {
        http
                .csrf(ServerHttpSecurity.CsrfSpec::disable)
                .exceptionHandling(exception -> exception
                        .authenticationEntryPoint((exchange, ex) ->
                                errorResponseWriter.write(
                                        exchange,
                                        HttpStatus.UNAUTHORIZED,
                                        ErrorCode.UNAUTHORIZED,
                                        ErrorKey.AUTHENTICATION_FAILED,
                                        "Authentication required"
                                )
                        )
                        .accessDeniedHandler((exchange, ex) ->
                                errorResponseWriter.write(
                                        exchange,
                                        HttpStatus.FORBIDDEN,
                                        ErrorCode.FORBIDDEN,
                                        ErrorKey.INSUFFICIENT_PERMISSION,
                                        "Access denied"
                                )
                        )
                )
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

                    if (securityProperties.jwtEnabled()) {
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
