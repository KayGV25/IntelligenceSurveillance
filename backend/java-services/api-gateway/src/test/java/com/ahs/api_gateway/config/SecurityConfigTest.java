package com.ahs.api_gateway.config;

import com.ahs.api_gateway.exception.ErrorResponseWriter;
import org.junit.jupiter.api.Test;
import org.springframework.security.config.web.server.ServerHttpSecurity;
import org.springframework.security.web.server.SecurityWebFilterChain;
import tools.jackson.databind.ObjectMapper;

import static org.assertj.core.api.Assertions.assertThat;

class SecurityConfigTest {

    @Test
    void buildsPermissiveSecurityFilterChainWhenJwtIsDisabled() {
        SecurityConfig config = new SecurityConfig(
                new SecurityProperties(false),
                new ErrorResponseWriter(new ObjectMapper())
        );

        SecurityWebFilterChain chain = config.securityWebFilterChain(ServerHttpSecurity.http());

        assertThat(chain).isNotNull();
    }
}
