package com.ahs.auth_service.auth;

import org.springframework.data.repository.reactive.ReactiveCrudRepository;
import reactor.core.publisher.Mono;

import java.util.UUID;

public interface RefreshTokenRepository extends ReactiveCrudRepository<RefreshTokenEntity, UUID> {

    Mono<RefreshTokenEntity> findByTokenHash(String tokenHash);
}
