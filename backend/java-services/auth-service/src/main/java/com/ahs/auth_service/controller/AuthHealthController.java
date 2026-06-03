package com.ahs.auth_service.controller;

import com.ahs.common.response.ApiResponse;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.time.Instant;
import java.util.Map;

@RestController
public class AuthHealthController {

    @GetMapping("/api/v1/auth/health")
    public ApiResponse<Map<String,Object>> health(){
        return ApiResponse.success(Map.of(
                "service", "auth-service",
                "status", "UP",
                "timestamp", Instant.now().toString()
        ));
    }
}
