package com.ahs.auth_service.audit;

import lombok.RequiredArgsConstructor;
import org.springframework.r2dbc.core.DatabaseClient;
import org.springframework.stereotype.Repository;
import reactor.core.publisher.Mono;

import java.util.UUID;

@Repository
@RequiredArgsConstructor
public class AuthAuditRepository {

    private final DatabaseClient databaseClient;

    public Mono<Void> save(
            UUID userId,
            String action,
            String ipAddress,
            String userAgent,
            String traceId,
            String requestId,
            String correlationId
    ) {
        return databaseClient.sql("""
                        INSERT INTO auth.audit_logs (
                            user_id,
                            action,
                            ip_address,
                            user_agent,
                            trace_id,
                            request_id,
                            correlation_id
                        )
                        VALUES (
                            :userId,
                            :action,
                            :ipAddress,
                            :userAgent,
                            :traceId,
                            :requestId,
                            :correlationId
                        )
                        """)
                .bind("userId", userId)
                .bind("action", action)
                .bind("ipAddress", ipAddress)
                .bind("userAgent", userAgent)
                .bind("traceId", traceId)
                .bind("requestId", requestId)
                .bind("correlationId", correlationId)
                .then();
    }

    public Mono<Void> saveAnonymous(
            String action,
            String ipAddress,
            String userAgent,
            String traceId,
            String requestId,
            String correlationId
    ) {
        return databaseClient.sql("""
                        INSERT INTO auth.audit_logs (
                            action,
                            ip_address,
                            user_agent,
                            trace_id,
                            request_id,
                            correlation_id
                        )
                        VALUES (
                            :action,
                            :ipAddress,
                            :userAgent,
                            :traceId,
                            :requestId,
                            :correlationId
                        )
                        """)
                .bind("action", action)
                .bind("ipAddress", ipAddress)
                .bind("userAgent", userAgent)
                .bind("traceId", traceId)
                .bind("requestId", requestId)
                .bind("correlationId", correlationId)
                .then();
    }
}
