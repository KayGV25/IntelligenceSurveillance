package com.ahs.auth_service.user;

import org.springframework.data.repository.reactive.ReactiveCrudRepository;
import reactor.core.publisher.Mono;

import java.util.UUID;

public interface UserRepository extends ReactiveCrudRepository<UserEntity, UUID> {
    Mono<UserEntity> findByEmail(String email);
    Mono<UserEntity> save(UserEntity user);
}
