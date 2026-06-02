package com.ahs.api_gateway.config;

import io.swagger.v3.oas.models.OpenAPI;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class OpenApiConfigTest {

    @Test
    void apiGatewayOpenApiContainsGatewayMetadata() {
        OpenAPI openAPI = new OpenApiConfig().apiGatewayOpenAPI();

        assertThat(openAPI.getInfo().getTitle()).isEqualTo("Advanced Home Surveillance API Gateway");
        assertThat(openAPI.getInfo().getVersion()).isEqualTo("0.0.1");
        assertThat(openAPI.getInfo().getContact().getEmail()).isEqualTo("khuongvudang25@gmail.com");
        assertThat(openAPI.getServers())
                .singleElement()
                .satisfies(server -> {
                    assertThat(server.getUrl()).isEqualTo("http://localhost:8080");
                    assertThat(server.getDescription()).isEqualTo("Local API Gateway");
                });
    }
}
