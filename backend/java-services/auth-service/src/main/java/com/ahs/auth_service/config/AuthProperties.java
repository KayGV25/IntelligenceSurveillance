package com.ahs.auth_service.config;

import org.springframework.boot.context.properties.ConfigurationProperties;

@ConfigurationProperties(prefix = "ahs.auth")
public record AuthProperties(
        String issuer,
        long accessTokenMinutes,
        long refreshTokenDays,
        String keyId,
        String privateKeyPath,
        String publicKeyPath
) {
}
