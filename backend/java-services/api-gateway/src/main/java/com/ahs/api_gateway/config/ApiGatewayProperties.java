package com.ahs.api_gateway.config;

import org.springframework.boot.context.properties.ConfigurationProperties;

@ConfigurationProperties(prefix = "ahs.services")
public record ApiGatewayProperties(
        String authServiceUrl,
        String cameraServiceUrl,
        String mapServiceUrl,
        String alertServiceUrl,
        String incidentServiceUrl,
        String trackingServiceUrl,
        String playbackServiceUrl
) {
}