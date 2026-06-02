package com.ahs.api_gateway.config;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.cloud.gateway.route.Route;
import org.springframework.cloud.gateway.route.RouteLocator;
import reactor.test.StepVerifier;

import java.net.URI;
import java.util.Map;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest(properties = {
        "ahs.services.auth-service-url=http://auth-service:8081",
        "ahs.services.camera-service-url=http://camera-service:8101",
        "ahs.services.map-service-url=http://map-service:8105",
        "ahs.services.alert-service-url=http://alert-service:8106",
        "ahs.services.incident-service-url=http://incident-service:8107",
        "ahs.services.tracking-service-url=http://tracking-service:8108",
        "ahs.services.playback-service-url=http://playback-service:8104"
})
class ApiGatewayRouteConfigTest {

    @Autowired
    private RouteLocator routeLocator;

    @Test
    void authRouteUsesAuthServiceUri() {
        StepVerifier.create(routeLocator.getRoutes()
                        .filter(route -> route.getId().equals("auth-service"))
                        .single()
                        .map(Route::getUri))
                .assertNext(uri -> assertThat(uri).isEqualTo(URI.create("http://auth-service:8081")))
                .verifyComplete();
    }

    @Test
    void allRoutesUseTheirConfiguredServiceUris() {
        Map<String, URI> expectedUris = Map.of(
                "auth-service", URI.create("http://auth-service:8081"),
                "camera-service", URI.create("http://camera-service:8101"),
                "map-service", URI.create("http://map-service:8105"),
                "alert-service", URI.create("http://alert-service:8106"),
                "incident-service", URI.create("http://incident-service:8107"),
                "tracking-service", URI.create("http://tracking-service:8108"),
                "playback-service", URI.create("http://playback-service:8104")
        );

        StepVerifier.create(routeLocator.getRoutes()
                        .collectMap(Route::getId, Route::getUri))
                .assertNext(routes -> assertThat(routes).containsExactlyInAnyOrderEntriesOf(expectedUris))
                .verifyComplete();
    }
}
