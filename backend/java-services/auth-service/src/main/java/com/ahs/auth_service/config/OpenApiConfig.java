package com.ahs.auth_service.config;

import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.oas.models.servers.Server;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.List;

@Configuration
public class OpenApiConfig {

    @Bean
    public OpenAPI authServiceOpenAPI() {
        return new OpenAPI()
                .info(new Info()
                        .title("Advanced Home Surveillance Auth Service")
                        .version("0.0.1")
                        .description("Authentication and authorization service for Advanced Home Surveillance backend.")
                        .contact(new Contact()
                                .name("AHS Backend Team")
                                .email("khuongvudang25@gmail.com")
                                .url("https://github.com/KayGV25/IntelligenceSurveillance")))
                .servers(List.of(
                        new Server()
                                .url("http://localhost:8081")
                                .description("Local Auth Service")
                ));
    }
}
