package com.ahs.api_gateway.config;

import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class ApiGatewayPropertiesTest {

    @Test
    void exposesConfiguredServiceUrls() {
        ApiGatewayProperties properties = new ApiGatewayProperties(
                "http://auth",
                "http://camera",
                "http://map",
                "http://alert",
                "http://incident",
                "http://tracking",
                "http://playback"
        );

        assertThat(properties.authServiceUrl()).isEqualTo("http://auth");
        assertThat(properties.cameraServiceUrl()).isEqualTo("http://camera");
        assertThat(properties.mapServiceUrl()).isEqualTo("http://map");
        assertThat(properties.alertServiceUrl()).isEqualTo("http://alert");
        assertThat(properties.incidentServiceUrl()).isEqualTo("http://incident");
        assertThat(properties.trackingServiceUrl()).isEqualTo("http://tracking");
        assertThat(properties.playbackServiceUrl()).isEqualTo("http://playback");
    }
}
