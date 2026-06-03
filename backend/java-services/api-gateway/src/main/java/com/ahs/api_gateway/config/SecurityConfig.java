package com.ahs.api_gateway.config;

import com.ahs.api_gateway.exception.ErrorResponseWriter;
import com.ahs.api_gateway.security.JwtAuthenticationConverter;
import com.ahs.common.error.ErrorCode;
import com.ahs.common.error.ErrorKey;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.security.config.web.server.ServerHttpSecurity;
import org.springframework.security.oauth2.server.resource.authentication.ReactiveJwtAuthenticationConverterAdapter;
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
                        exchange
                                .pathMatchers("/api/v1/auth/logout").authenticated()
                                .pathMatchers("/api/v1/auth/me").authenticated()

                                .pathMatchers(HttpMethod.GET, "/api/v1/cameras/**")
                                .hasAnyAuthority("camera.read", "ROLE_SUPER_ADMIN", "ROLE_ADMIN", "ROLE_SECURITY_OPERATOR", "ROLE_VIEWER")

                                .pathMatchers(HttpMethod.POST, "/api/v1/cameras/**")
                                .hasAnyAuthority("camera.write", "ROLE_SUPER_ADMIN", "ROLE_ADMIN")

                                .pathMatchers(HttpMethod.PUT, "/api/v1/cameras/**")
                                .hasAnyAuthority("camera.write", "ROLE_SUPER_ADMIN", "ROLE_ADMIN")

                                .pathMatchers(HttpMethod.DELETE, "/api/v1/cameras/**")
                                .hasAnyAuthority("camera.write", "ROLE_SUPER_ADMIN", "ROLE_ADMIN")

                                .pathMatchers(HttpMethod.GET, "/api/v1/maps/**")
                                .hasAnyAuthority("map.read", "ROLE_SUPER_ADMIN", "ROLE_ADMIN", "ROLE_SECURITY_OPERATOR", "ROLE_VIEWER")

                                .pathMatchers("/api/v1/alerts/**")
                                .hasAnyAuthority("alert.read", "alert.write", "ROLE_SUPER_ADMIN", "ROLE_ADMIN", "ROLE_SECURITY_OPERATOR")

                                .pathMatchers("/api/v1/incidents/**")
                                .hasAnyAuthority("incident.read", "incident.write", "ROLE_SUPER_ADMIN", "ROLE_ADMIN", "ROLE_SECURITY_OPERATOR", "ROLE_INVESTIGATOR")

                                .pathMatchers("/api/v1/playback/**")
                                .hasAnyAuthority("playback.read", "ROLE_SUPER_ADMIN", "ROLE_ADMIN", "ROLE_SECURITY_OPERATOR", "ROLE_INVESTIGATOR")

                                .pathMatchers("/api/v1/tracking/**")
                                .hasAnyAuthority("tracking.read", "ROLE_SUPER_ADMIN", "ROLE_ADMIN", "ROLE_SECURITY_OPERATOR", "ROLE_INVESTIGATOR")

                                .anyExchange()
                                .authenticated();
                    } else {
                        exchange.anyExchange().permitAll();
                    }
                });

        if (securityProperties.jwtEnabled()) {
            http.oauth2ResourceServer(oauth2 -> oauth2
                    .jwt(jwt -> jwt.jwtAuthenticationConverter(
                            new ReactiveJwtAuthenticationConverterAdapter(
                                    new JwtAuthenticationConverter()
                            )
                    ))
            );
        }
        return http.build();
    }
}
