package com.ahs.auth_service.auth;

import com.ahs.auth_service.audit.AuthAuditRepository;
import com.ahs.auth_service.audit.RequestAuditContext;
import com.ahs.auth_service.auth.dto.*;
import com.ahs.auth_service.user.UserAuthorizationRepository;
import com.ahs.auth_service.user.UserEntity;
import com.ahs.auth_service.user.UserRepository;
import com.ahs.common.error.BusinessException;
import com.ahs.common.error.ErrorCode;
import com.ahs.common.error.ErrorKey;
import lombok.RequiredArgsConstructor;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;

import java.time.OffsetDateTime;

@Service
@RequiredArgsConstructor
public class AuthApplicationService {

    private final UserRepository userRepository;
    private final UserAuthorizationRepository authorizationRepository;
    private final PasswordEncoder passwordEncoder;
    private final JwtService jwtService;
    private final RefreshTokenService refreshTokenService;
    private final AuthAuditRepository auditRepository;

    public Mono<AuthResponse> login(
            LoginRequest request,
            RequestAuditContext auditContext
    ) {
        return userRepository.findByEmail(request.email())
                .switchIfEmpty(Mono.defer(() ->
                        auditRepository.saveAnonymous(
                                "LOGIN_FAILED",
                                auditContext.ipAddress(),
                                auditContext.userAgent(),
                                auditContext.traceId(),
                                auditContext.requestId(),
                                auditContext.correlationId()
                        )
                                .then(Mono.error(new BusinessException(
                                        ErrorCode.UNAUTHORIZED,
                                        ErrorKey.INVALID_CREDENTIALS,
                                        "Invalid email or password"
                                )))
                ))
                .flatMap(user -> validatePassword(user, request.password())
                        .flatMap(validUser -> issueTokens(validUser)
                                .flatMap(response ->
                                        auditRepository.save(
                                                        validUser.getId(),
                                                        "LOGIN_SUCCESS",
                                                        auditContext.ipAddress(),
                                                        auditContext.userAgent(),
                                                        auditContext.traceId(),
                                                        auditContext.requestId(),
                                                        auditContext.correlationId()
                                                )
                                                .thenReturn(response)
                                )
                        )
                        .onErrorResume(BusinessException.class, ex ->
                                auditRepository.save(
                                                user.getId(),
                                                "LOGIN_FAILED",
                                                auditContext.ipAddress(),
                                                auditContext.userAgent(),
                                                auditContext.traceId(),
                                                auditContext.requestId(),
                                                auditContext.correlationId()
                                        )
                                        .then(Mono.error(ex))
                        )
                );
    }

    public Mono<AuthResponse> refresh(
            RefreshTokenRequest request,
            RequestAuditContext auditContext
    ) {
        return refreshTokenService.validateRefreshToken(request.refreshToken())
                .flatMap(token -> userRepository.findById(token.getUserId())
                        .switchIfEmpty(Mono.error(new BusinessException(
                                ErrorCode.UNAUTHORIZED,
                                ErrorKey.INVALID_REFRESH_TOKEN,
                                "Invalid refresh token"
                        )))
                        .flatMap(this::issueTokens)
                        .flatMap(response ->
                                refreshTokenService.revokeRefreshToken(request.refreshToken())
                                        .then(auditRepository.save(
                                                token.getUserId(),
                                                "TOKEN_REFRESHED",
                                                auditContext.ipAddress(),
                                                auditContext.userAgent(),
                                                auditContext.traceId(),
                                                auditContext.requestId(),
                                                auditContext.correlationId()
                                        ))
                                        .thenReturn(response)
                        )
                );
    }

    public Mono<MeResponse> me(String authorizationHeader) {
        String token = extractBearerToken(authorizationHeader);
        return Mono.fromCallable(() -> jwtService.validateAndExtractUser(token));
    }

    public Mono<Void> logout(
            LogoutRequest request,
            RequestAuditContext auditContext
    ) {
        return refreshTokenService.validateRefreshToken(request.refreshToken())
                .flatMap(token ->
                        refreshTokenService.revokeRefreshToken(request.refreshToken())
                                .then(auditRepository.save(
                                        token.getUserId(),
                                        "LOGOUT_SUCCESS",
                                        auditContext.ipAddress(),
                                        auditContext.userAgent(),
                                        auditContext.traceId(),
                                        auditContext.requestId(),
                                        auditContext.correlationId()
                                ))
                );
    }

    private Mono<AuthResponse> issueTokens(UserEntity user) {
        return Mono.zip(
                authorizationRepository.findRoleNamesByUserId(user.getId()),
                authorizationRepository.findPermissionNamesByUserId(user.getId())
        ).flatMap(tuple -> {
            var roles = tuple.getT1();
            var permissions = tuple.getT2();

            String accessToken = jwtService.generateAccessToken(
                    user,
                    roles,
                    permissions
            );

            return refreshTokenService.createRefreshToken(user.getId())
                    .map(refreshToken -> new AuthResponse(
                            accessToken,
                            refreshToken,
                            jwtService.accessTokenExpiresInSeconds(),
                            "Bearer",
                            new AuthResponse.UserProfile(
                                    user.getId().toString(),
                                    user.getEmail(),
                                    roles,
                                    permissions
                            )
                    ));
        });
    }


    private Mono<UserEntity> validatePassword(UserEntity user, String rawPassword) {
        if (!Boolean.TRUE.equals(user.getIsEnabled())) {
            return Mono.error(new BusinessException(
                    ErrorCode.FORBIDDEN,
                    ErrorKey.USER_DISABLED,
                    "User account is disabled"
            ));
        }

        if (!passwordEncoder.matches(rawPassword, user.getPasswordHash())) {
            int attempts = user.getFailedLoginAttempts() == null
                    ? 1
                    : user.getFailedLoginAttempts() + 1;

            user.setFailedLoginAttempts(attempts);

            if (attempts >= 5) {
                user.setIsAccountLocked(true);
                user.setLockedUntil(OffsetDateTime.now().plusMinutes(15));
            }

            return userRepository.save(user)
                    .then(Mono.error(new BusinessException(
                            ErrorCode.UNAUTHORIZED,
                            ErrorKey.INVALID_CREDENTIALS,
                            "Invalid email or password"
                    )));
        }

        if (!passwordEncoder.matches(rawPassword, user.getPasswordHash())) {
            int attempts = user.getFailedLoginAttempts() == null
                    ? 1
                    : user.getFailedLoginAttempts() + 1;

            user.setFailedLoginAttempts(attempts);

            if (attempts >= 5) {
                user.setIsAccountLocked(true);
                user.setLockedUntil(OffsetDateTime.now().plusMinutes(15));
            }

            return userRepository.save(user)
                    .then(Mono.error(new BusinessException(
                            ErrorCode.UNAUTHORIZED,
                            ErrorKey.INVALID_CREDENTIALS,
                            "Invalid email or password"
                    )));
        }

        user.setFailedLoginAttempts(0);
        user.setIsAccountLocked(false);
        user.setLockedUntil(null);
        return userRepository.save(user).then(Mono.just(user));
    }

    private String extractBearerToken(String authorizationHeader) {
        if (authorizationHeader == null ||
                !authorizationHeader.startsWith("Bearer ")) {
            throw new BusinessException(
                    ErrorCode.UNAUTHORIZED,
                    ErrorKey.ACCESS_TOKEN_INVALID,
                    "Missing or invalid Authorization header"
            );
        }

        return authorizationHeader.substring("Bearer ".length());
    }
}
