package com.ahs.api_gateway.config;

import org.junit.jupiter.api.Test;
import org.springframework.security.config.web.server.ServerHttpSecurity;
import org.springframework.security.web.server.SecurityWebFilterChain;

import static org.assertj.core.api.Assertions.assertThat;

class SecurityConfigTest {

    @Test
    void buildsPermissiveSecurityFilterChainWhenJwtIsDisabled() {
        SecurityConfig config = new SecurityConfig(new SecurityProperties(false));

        SecurityWebFilterChain chain = config.securityWebFilterChain(ServerHttpSecurity.http());

        assertThat(chain).isNotNull();
    }
}
