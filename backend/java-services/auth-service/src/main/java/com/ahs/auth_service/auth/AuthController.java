package com.ahs.auth_service.auth;

import com.ahs.auth_service.audit.RequestAuditContext;
import com.ahs.auth_service.auth.dto.*;
import com.ahs.common.response.ApiResponse;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.server.reactive.ServerHttpRequest;
import org.springframework.web.bind.annotation.*;
import reactor.core.publisher.Mono;

@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/auth")
public class AuthController {

    private final AuthApplicationService authService;

    @PostMapping("/login")
    public Mono<ApiResponse<AuthResponse>> login(
            @Valid @RequestBody LoginRequest request,
            ServerHttpRequest httpRequest
    ) {
        return authService.login(
                request,
                RequestAuditContext.from(httpRequest)
        ).map(ApiResponse::success);
    }

    @PostMapping("/refresh")
    public Mono<ApiResponse<AuthResponse>> refresh(
            @Valid @RequestBody RefreshTokenRequest request,
            ServerHttpRequest httpRequest
    ) {
        return authService.refresh(
                request,
                RequestAuditContext.from(httpRequest)
        ).map(ApiResponse::success);
    }

    @PostMapping("/logout")
    public Mono<ApiResponse<String>> logout(
            @Valid @RequestBody LogoutRequest request,
            ServerHttpRequest httpRequest
    ) {
        return authService.logout(
                request,
                RequestAuditContext.from(httpRequest)
        ).thenReturn(ApiResponse.success("Logged out"));
    }

    @GetMapping("/me")
    public Mono<ApiResponse<MeResponse>> me(
            @RequestHeader(value = "Authorization", required = false)
            String authorizationHeader
    ) {
        return authService.me(authorizationHeader)
                .map(ApiResponse::success);
    }
}
