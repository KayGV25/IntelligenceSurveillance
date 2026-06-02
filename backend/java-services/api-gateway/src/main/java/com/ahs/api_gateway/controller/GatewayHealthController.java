package com.ahs.api_gateway.controller;

import com.ahs.common.response.ApiResponse;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.time.Instant;
import java.util.Map;

@RestController
public class GatewayHealthController {

    @GetMapping("/api/v1/gateway/health")
    public ApiResponse<Map<String, Object>> health() {
        return ApiResponse.success(
                Map.of(
                        "service", "api-gateway",
                        "status", "UP",
                        "timestamp", Instant.now().toString())
        );
    }
}
