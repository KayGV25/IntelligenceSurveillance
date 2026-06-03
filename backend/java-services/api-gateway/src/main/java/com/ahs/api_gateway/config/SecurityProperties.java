package com.ahs.api_gateway.config;

import org.springframework.boot.context.properties.ConfigurationProperties;

@ConfigurationProperties(prefix = "ahs.security")
public record SecurityProperties(
        boolean jwtEnabled,
        String jwtIssuer,
        String jwkSetUri
) {
}
