INSERT INTO auth.users (
    email,
    password_hash,
    first_name,
    last_name,
    is_enabled,
    is_account_locked
)
VALUES (
           'admin@ahs.local',
           '$2a$10$cobbIvIxVGs6/Mm0jXWdfupK78JvHzaLffMHAIcO4hFcy1uzYc08G',
           'System',
           'Admin',
           TRUE,
           FALSE
       )
ON CONFLICT (email)
    DO UPDATE SET
                  password_hash = EXCLUDED.password_hash,
                  first_name = EXCLUDED.first_name,
                  last_name = EXCLUDED.last_name,
                  is_enabled = EXCLUDED.is_enabled,
                  is_account_locked = EXCLUDED.is_account_locked,
                  updated_at = NOW();

INSERT INTO auth.user_roles (user_id, role_id)
SELECT u.id, r.id
FROM auth.users u
         JOIN auth.roles r ON r.name = 'SUPER_ADMIN'
WHERE u.email = 'admin@ahs.local'
ON CONFLICT DO NOTHING;