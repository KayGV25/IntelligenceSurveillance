package com.ahs.auth_service.auth;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.Id;
import org.springframework.data.relational.core.mapping.Table;

import java.time.OffsetDateTime;
import java.util.UUID;

@Table(schema = "auth", name = "refresh_tokens")
public class RefreshTokenEntity {

    @Id
    @Getter
    @Setter
    private UUID id;

    @Getter
    @Setter
    private UUID userId;

    @Getter
    @Setter
    private String tokenHash;

    @Getter
    @Setter
    private Boolean revoked;

    @Getter
    @Setter
    private OffsetDateTime expiresAt;

    @Getter
    @Setter
    private OffsetDateTime createdAt;
}
