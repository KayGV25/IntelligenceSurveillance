INSERT INTO auth.roles (name, description)
VALUES
    ('SUPER_ADMIN', 'Full system access'),
    ('ADMIN', 'Organization administrator'),
    ('SECURITY_OPERATOR', 'Security monitoring operator'),
    ('INVESTIGATOR', 'Historical footage and incident investigator'),
    ('VIEWER', 'Read-only viewer')
    ON CONFLICT (name) DO NOTHING;

INSERT INTO auth.permissions (name, description)
VALUES
    ('user.read', 'View users'),
    ('user.write', 'Create and update users'),
    ('camera.read', 'View cameras'),
    ('camera.write', 'Create and update cameras'),
    ('map.read', 'View maps'),
    ('map.write', 'Create and update maps'),
    ('alert.read', 'View alerts'),
    ('alert.write', 'Update alerts'),
    ('incident.read', 'View incidents'),
    ('incident.write', 'Create and update incidents'),
    ('playback.read', 'View playback'),
    ('tracking.read', 'View tracking data')
    ON CONFLICT (name) DO NOTHING;

INSERT INTO auth.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM auth.roles r
         CROSS JOIN auth.permissions p
WHERE r.name = 'SUPER_ADMIN'
    ON CONFLICT DO NOTHING;

INSERT INTO auth.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM auth.roles r
         JOIN auth.permissions p ON p.name IN (
                                               'user.read',
                                               'camera.read',
                                               'camera.write',
                                               'map.read',
                                               'map.write',
                                               'alert.read',
                                               'alert.write',
                                               'incident.read',
                                               'incident.write',
                                               'playback.read',
                                               'tracking.read'
    )
WHERE r.name = 'ADMIN'
    ON CONFLICT DO NOTHING;

INSERT INTO auth.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM auth.roles r
         JOIN auth.permissions p ON p.name IN (
                                               'camera.read',
                                               'map.read',
                                               'alert.read',
                                               'alert.write',
                                               'incident.read',
                                               'incident.write',
                                               'playback.read',
                                               'tracking.read'
    )
WHERE r.name = 'SECURITY_OPERATOR'
    ON CONFLICT DO NOTHING;

INSERT INTO auth.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM auth.roles r
         JOIN auth.permissions p ON p.name IN (
                                               'incident.read',
                                               'incident.write',
                                               'playback.read',
                                               'tracking.read'
    )
WHERE r.name = 'INVESTIGATOR'
    ON CONFLICT DO NOTHING;

INSERT INTO auth.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM auth.roles r
         JOIN auth.permissions p ON p.name IN (
                                               'camera.read',
                                               'map.read',
                                               'alert.read'
    )
WHERE r.name = 'VIEWER'
    ON CONFLICT DO NOTHING;