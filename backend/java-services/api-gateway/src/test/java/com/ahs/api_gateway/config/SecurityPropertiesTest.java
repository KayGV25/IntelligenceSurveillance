package com.ahs.api_gateway.config;

import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class SecurityPropertiesTest {

    @Test
    void exposesJwtEnabledFlag() {
        assertThat(new SecurityProperties(true).jwtEnabled()).isTrue();
        assertThat(new SecurityProperties(false).jwtEnabled()).isFalse();
    }
}
