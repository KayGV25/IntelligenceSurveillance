package com.ahs.auth_service.auth;

import com.ahs.auth_service.config.AuthProperties;
import com.ahs.common.error.BusinessException;
import com.ahs.common.error.ErrorCode;
import com.ahs.common.error.ErrorKey;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;

import java.security.SecureRandom;
import java.time.OffsetDateTime;
import java.util.Base64;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class RefreshTokenService {

    private final SecureRandom secureRandom = new SecureRandom();

    private final RefreshTokenRepository refreshTokenRepository;
    private final TokenHashService tokenHashService;
    private final AuthProperties  authProperties;

    public Mono<String> createRefreshToken(UUID userId){
        String rawToken = generateSecureToken();
        String tokenHash = tokenHashService.sha256(rawToken);

        RefreshTokenEntity entity = new RefreshTokenEntity();
        entity.setUserId(userId);
        entity.setTokenHash(tokenHash);
        entity.setRevoked(false);
        entity.setExpiresAt(OffsetDateTime.now().plusDays(authProperties.refreshTokenDays()));

        return refreshTokenRepository.save(entity)
                .thenReturn(rawToken);
    }

    public Mono<RefreshTokenEntity> validateRefreshToken(String rawToken){
        String tokenHash = tokenHashService.sha256(rawToken);

        return refreshTokenRepository.findByTokenHash(tokenHash)
                .switchIfEmpty(Mono.error(new BusinessException(
                        ErrorCode.UNAUTHORIZED,
                        ErrorKey.INVALID_REFRESH_TOKEN,
                        "Invalid refresh token"
                )))
                .flatMap(token -> {
                    if (Boolean.TRUE.equals(token.getRevoked())) {
                        return Mono.error(new BusinessException(
                                ErrorCode.UNAUTHORIZED,
                                ErrorKey.REFRESH_TOKEN_REVOKED,
                                "Refresh token has been revoked"
                        ));
                    }

                    if (token.getExpiresAt().isBefore(OffsetDateTime.now())) {
                        return Mono.error(new BusinessException(
                                ErrorCode.UNAUTHORIZED,
                                ErrorKey.REFRESH_TOKEN_EXPIRED,
                                "Refresh token has expired"
                        ));
                    }

                    return Mono.just(token);
                });
    }

    public Mono<Void> revokeRefreshToken(String rawToken){
        String tokenHash = tokenHashService.sha256(rawToken);

        return refreshTokenRepository.findByTokenHash(tokenHash)
                .flatMap(token -> {
                    token.setRevoked(true);
                    return refreshTokenRepository.save(token);
                })
                .then();
    }

    private String generateSecureToken() {
        byte[] bytes = new byte[64];
        secureRandom.nextBytes(bytes);
        return Base64.getUrlEncoder().withoutPadding().encodeToString(bytes);
    }
}
