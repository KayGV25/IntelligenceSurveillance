FROM eclipse-temurin:21-jdk-jammy AS build

WORKDIR /workspace

COPY gradlew .
COPY gradle gradle
COPY settings.gradle.kts .
COPY build.gradle.kts .
COPY api-gateway/build.gradle.kts api-gateway/build.gradle.kts

RUN chmod +x gradlew
RUN ./gradlew :api-gateway:dependencies --no-daemon

COPY api-gateway/src api-gateway/src

RUN ./gradlew :api-gateway:bootJar --no-daemon

FROM eclipse-temurin:21-jre-jammy

WORKDIR /app

COPY --from=build /workspace/api-gateway/build/libs/*.jar app.jar

EXPOSE 8080

ENTRYPOINT ["java", "-jar", "app.jar"]