package com.ahs.api_gateway.filter;

import com.ahs.common.security.UserContextHeaders;
import org.springframework.cloud.gateway.filter.GatewayFilterChain;
import org.springframework.cloud.gateway.filter.GlobalFilter;
import org.springframework.core.Ordered;
import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.stereotype.Component;
import org.springframework.web.server.ServerWebExchange;
import reactor.core.publisher.Mono;

import java.util.List;

@Component
public class UserContextPropagationFilter implements GlobalFilter, Ordered {

    @Override
    public Mono<Void> filter(
            ServerWebExchange exchange,
            GatewayFilterChain chain
    ) {
        return exchange.getPrincipal()
                .cast(Authentication.class)
                .flatMap(authentication -> {
                    if (!(authentication.getPrincipal() instanceof Jwt jwt)) {
                        return chain.filter(exchange);
                    }

                    String userId = jwt.getSubject();
                    String email = jwt.getClaimAsString("email");
                    List<String> roles = jwt.getClaimAsStringList("roles");
                    List<String> permissions = jwt.getClaimAsStringList("permissions");

                    var mutatedRequest = exchange.getRequest()
                            .mutate()
                            .header(UserContextHeaders.USER_ID, userId == null ? "" : userId)
                            .header(UserContextHeaders.USER_EMAIL, email == null ? "" : email)
                            .header(UserContextHeaders.USER_ROLES, roles == null ? "" : String.join(",", roles))
                            .header(UserContextHeaders.USER_PERMISSIONS, permissions == null ? "" : String.join(",", permissions))
                            .build();

                    return chain.filter(exchange.mutate().request(mutatedRequest).build());
                })
                .switchIfEmpty(chain.filter(exchange));
    }

    @Override
    public int getOrder() {
        return -50;
    }
}
