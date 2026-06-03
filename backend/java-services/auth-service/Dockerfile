FROM eclipse-temurin:21-jdk-jammy AS build

WORKDIR /workspace

COPY gradlew .
COPY gradle gradle
COPY settings.gradle.kts .
COPY build.gradle.kts .

COPY libs/common-response/build.gradle.kts libs/common-response/build.gradle.kts
COPY libs/common-error/build.gradle.kts libs/common-error/build.gradle.kts
COPY libs/common-security/build.gradle.kts libs/common-security/build.gradle.kts
COPY libs/common-observability/build.gradle.kts libs/common-observability/build.gradle.kts

COPY api-gateway/build.gradle.kts api-gateway/build.gradle.kts
COPY auth-service/build.gradle.kts auth-service/build.gradle.kts

RUN chmod +x gradlew
RUN ./gradlew :auth-service:dependencies --no-daemon

COPY libs libs
COPY api-gateway/src api-gateway/src
COPY auth-service/src auth-service/src

RUN ./gradlew :auth-service:bootJar --no-daemon

FROM eclipse-temurin:21-jre-jammy

WORKDIR /app

COPY --from=build /workspace/auth-service/build/libs/*.jar app.jar

EXPOSE 8081

ENTRYPOINT ["java", "-jar", "app.jar"]