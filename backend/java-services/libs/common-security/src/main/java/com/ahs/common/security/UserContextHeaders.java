package com.ahs.common.security;

public final class UserContextHeaders {

    private UserContextHeaders() {}

    public static final String USER_ID = "X-User-Id";
    public static final String USER_EMAIL = "X-User-Email";
    public static final String USER_ROLES = "X-User-Roles";
    public static final String USER_PERMISSIONS = "X-User-Permissions";
}
