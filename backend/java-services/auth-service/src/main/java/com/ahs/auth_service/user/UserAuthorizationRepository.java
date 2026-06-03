package com.ahs.auth_service.user;

import lombok.RequiredArgsConstructor;
import org.springframework.r2dbc.core.DatabaseClient;
import org.springframework.stereotype.Repository;
import reactor.core.publisher.Mono;

import java.util.List;
import java.util.Objects;
import java.util.UUID;

@Repository
@RequiredArgsConstructor
public class UserAuthorizationRepository {

    private final DatabaseClient databaseClient;

    public Mono<List<String>> findRoleNamesByUserId(UUID userId) {
        return databaseClient.sql("""
                        SELECT r.name
                        FROM auth.roles r
                        JOIN auth.user_roles ur ON ur.role_id = r.id
                        WHERE ur.user_id = :userId
                        """)
                .bind("userId", userId)
                .map((row, metadata) -> Objects.requireNonNull(row.get("name", String.class)))
                .all()
                .collectList();
    }

    public Mono<List<String>> findPermissionNamesByUserId(UUID userId) {
        return databaseClient.sql("""
                        SELECT DISTINCT p.name
                        FROM auth.permissions p
                        JOIN auth.role_permissions rp ON rp.permission_id = p.id
                        JOIN auth.user_roles ur ON ur.role_id = rp.role_id
                        WHERE ur.user_id = :userId
                        """)
                .bind("userId", userId)
                .map((row, metadata) -> Objects.requireNonNull(row.get("name", String.class)))
                .all()
                .collectList();
    }
}
