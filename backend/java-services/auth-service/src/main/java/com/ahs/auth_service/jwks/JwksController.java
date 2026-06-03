package com.ahs.auth_service.jwks;

import com.ahs.auth_service.auth.RsaKeyProvider;
import com.ahs.auth_service.config.AuthProperties;
import com.nimbusds.jose.JWSAlgorithm;
import com.nimbusds.jose.jwk.JWKSet;
import com.nimbusds.jose.jwk.RSAKey;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.security.interfaces.RSAPublicKey;
import java.util.Map;

@RestController
@RequiredArgsConstructor
public class JwksController {

    private final AuthProperties properties;
    private final RsaKeyProvider rsaKeyProvider;

    @GetMapping("/.well-known/jwks.json")
    public Map<String, Object> jwks() {
        RSAPublicKey publicKey =
                rsaKeyProvider.loadPublicKey(properties.publicKeyPath());

        RSAKey rsaKey = new RSAKey.Builder(publicKey)
                .keyID(properties.keyId())
                .algorithm(JWSAlgorithm.RS256)
                .build();

        return new JWKSet(rsaKey).toJSONObject();
    }
}
