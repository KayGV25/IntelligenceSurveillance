pluginManagement {
    repositories {
        mavenCentral()
        gradlePluginPortal()
    }
}

dependencyResolutionManagement {
    repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS)

    repositories {
        mavenCentral()
    }
}

rootProject.name = "ahs-java-services"

include(
    "api-gateway",
    "auth-service",
    "libs:common-response",
    "libs:common-error",
    "libs:common-security",
    "libs:common-observability"
)