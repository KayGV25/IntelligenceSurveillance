package com.ahs.auth_service.user;

import lombok.Getter;
import lombok.Setter;
import org.springframework.data.annotation.Id;
import org.springframework.data.relational.core.mapping.Table;

import java.time.OffsetDateTime;
import java.util.UUID;

@Table(schema = "auth", name = "users")
public class UserEntity {

    @Getter
    @Setter
    @Id
    private UUID id;

    @Getter
    private String email;

    @Getter
    private String passwordHash;

    private String firstName;

    private String lastName;

    @Getter
    private Boolean isEnabled;

    @Getter
    @Setter
    private Boolean isAccountLocked;

    @Getter
    @Setter
    private Integer failedLoginAttempts;

    @Getter
    @Setter
    private OffsetDateTime lockedUntil;

    private OffsetDateTime createdAt;

    private OffsetDateTime updatedAt;

}
