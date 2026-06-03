package com.ahs.auth_service.auth.dto;

import java.util.List;

public record AuthResponse(
        String accessToken,
        String refreshToken,
        long expiresInSeconds,
        String tokenType,
        UserProfile user
) {
    public record UserProfile(
            String id,
            String email,
            List<String> roles,
            List<String> permissions
    ) {}
}
