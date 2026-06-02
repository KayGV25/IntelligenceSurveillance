package com.ahs.api_gateway.controller;

import com.ahs.common.response.ApiResponse;
import org.junit.jupiter.api.Test;

import java.time.Instant;
import java.util.Map;

import static org.assertj.core.api.Assertions.assertThat;

class GatewayHealthControllerTest {

    private final GatewayHealthController controller = new GatewayHealthController();

    @Test
    void healthReturnsGatewayStatusPayload() {
        ApiResponse<Map<String, Object>> response = controller.health();

        assertThat(response.success()).isTrue();
        assertThat(response.error()).isNull();
        assertThat(response.data())
                .containsEntry("service", "api-gateway")
                .containsEntry("status", "UP");
        assertThat(Instant.parse(response.data().get("timestamp").toString())).isNotNull();
        assertThat(response.timestamp()).isNotNull();
    }
}
