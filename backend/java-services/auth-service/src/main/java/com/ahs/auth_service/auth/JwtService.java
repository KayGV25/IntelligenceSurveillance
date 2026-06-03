package com.ahs.auth_service.auth;

import com.ahs.auth_service.auth.dto.MeResponse;
import com.ahs.auth_service.config.AuthProperties;
import com.ahs.auth_service.user.UserEntity;
import com.ahs.common.error.BusinessException;
import com.ahs.common.error.ErrorCode;
import com.ahs.common.error.ErrorKey;
import com.nimbusds.jose.JOSEObjectType;
import com.nimbusds.jose.JWSAlgorithm;
import com.nimbusds.jose.JWSHeader;
import com.nimbusds.jose.crypto.RSASSASigner;
import com.nimbusds.jose.crypto.RSASSAVerifier;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.SignedJWT;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

import java.security.interfaces.RSAPrivateKey;
import java.security.interfaces.RSAPublicKey;
import java.time.Instant;
import java.util.Date;
import java.util.List;

@Service
@RequiredArgsConstructor
public class JwtService {
    private final AuthProperties properties;
    private final RsaKeyProvider rsaKeyProvider;

    public String generateAccessToken(
            UserEntity user,
            List<String> roles,
            List<String> permissions
    ) {
        try {
            Instant now = Instant.now();
            Instant expiresAt = now.plusSeconds(properties.accessTokenMinutes() * 60);

            JWTClaimsSet claims = new JWTClaimsSet.Builder()
                    .issuer(properties.issuer())
                    .subject(user.getId().toString())
                    .claim("email", user.getEmail())
                    .claim("roles", roles)
                    .claim("permissions", permissions)
                    .issueTime(Date.from(now))
                    .expirationTime(Date.from(expiresAt))
                    .build();

            JWSHeader header = new JWSHeader.Builder(JWSAlgorithm.RS256)
                    .type(JOSEObjectType.JWT)
                    .keyID(properties.keyId())
                    .build();

            SignedJWT jwt = new SignedJWT(
                    header,
                    claims
            );

            RSAPrivateKey privateKey =
                    rsaKeyProvider.loadPrivateKey(properties.privateKeyPath());

            jwt.sign(new RSASSASigner(privateKey));

            return jwt.serialize();
        } catch (Exception e) {
            throw new IllegalStateException("Failed to generate access token", e);
        }
    }

    public MeResponse validateAndExtractUser(String token) {
        try {
            SignedJWT signedJWT = SignedJWT.parse(token);

            RSAPublicKey publicKey =
                    rsaKeyProvider.loadPublicKey(properties.publicKeyPath());

            boolean valid = signedJWT.verify(new RSASSAVerifier(publicKey));

            if (!valid) {
                throw invalidToken();
            }

            JWTClaimsSet claims = signedJWT.getJWTClaimsSet();

            if (!properties.issuer().equals(claims.getIssuer())) {
                throw invalidToken();
            }

            if (claims.getExpirationTime() == null ||
                    claims.getExpirationTime().before(new Date())) {
                throw new BusinessException(
                        ErrorCode.UNAUTHORIZED,
                        ErrorKey.ACCESS_TOKEN_EXPIRED,
                        "Access token has expired"
                );
            }

            List<String> roles = claims.getStringListClaim("roles");
            List<String> permissions = claims.getStringListClaim("permissions");

            return new MeResponse(
                    claims.getSubject(),
                    claims.getStringClaim("email"),
                    roles == null ? List.of() : roles,
                    permissions == null ? List.of() : permissions
            );
        } catch (BusinessException e) {
            throw e;
        } catch (Exception e) {
            throw invalidToken();
        }
    }

    public long accessTokenExpiresInSeconds() {
        return properties.accessTokenMinutes() * 60;
    }

    private BusinessException invalidToken() {
        return new BusinessException(
                ErrorCode.UNAUTHORIZED,
                ErrorKey.ACCESS_TOKEN_INVALID,
                "Invalid access token"
        );
    }
}
