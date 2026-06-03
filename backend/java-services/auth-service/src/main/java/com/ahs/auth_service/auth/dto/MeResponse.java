package com.ahs.auth_service.auth.dto;

import java.util.List;

public record MeResponse(
        String id,
        String email,
        List<String> roles,
        List<String> permissions
) {
}
