-- ============================================================================
-- BASIC CRUD OPERATIONS
-- ============================================================================

-- name: CreateOneRolePermission :one
INSERT INTO role_permissions (id, role_id, permission_id)
VALUES ($1, $2, $3)
ON CONFLICT (role_id, permission_id) DO NOTHING
RETURNING *;

-- name: GetRolePermissionById :one
SELECT rp.*,
       json_build_object(
           'id', r.id,
           'title_ru', r.title_ru,
           'title_en', r.title_en,
           'title_kk', r.title_kk,
           'value', r.value,
           'description_ru', r.description_ru,
           'description_en', r.description_en,
           'description_kk', r.description_kk,
           'created_at', r.created_at,
           'updated_at', r.updated_at,
           'deleted_at', r.deleted_at
       ) as role,
       json_build_object(
           'id', p.id,
           'title_ru', p.title_ru,
           'title_en', p.title_en,
           'title_kk', p.title_kk,
           'value', p.value,
           'description_ru', p.description_ru,
           'description_en', p.description_en,
           'description_kk', p.description_kk,
           'created_at', p.created_at,
           'updated_at', p.updated_at,
           'deleted_at', p.deleted_at
       ) as permission
FROM role_permissions rp
INNER JOIN roles r ON rp.role_id = r.id
INNER JOIN permissions p ON rp.permission_id = p.id
WHERE rp.id = $1
  AND r.deleted_at IS NULL
  AND p.deleted_at IS NULL;

-- name: GetRolePermissionByRoleAndPermission :one
SELECT rp.*,
       json_build_object(
           'id', r.id,
           'title_ru', r.title_ru,
           'title_en', r.title_en,
           'title_kk', r.title_kk,
           'value', r.value,
           'description_ru', r.description_ru,
           'description_en', r.description_en,
           'description_kk', r.description_kk,
           'created_at', r.created_at,
           'updated_at', r.updated_at,
           'deleted_at', r.deleted_at
       ) as role,
       json_build_object(
           'id', p.id,
           'title_ru', p.title_ru,
           'title_en', p.title_en,
           'title_kk', p.title_kk,
           'value', p.value,
           'description_ru', p.description_ru,
           'description_en', p.description_en,
           'description_kk', p.description_kk,
           'created_at', p.created_at,
           'updated_at', p.updated_at,
           'deleted_at', p.deleted_at
       ) as permission
FROM role_permissions rp
INNER JOIN roles r ON rp.role_id = r.id
INNER JOIN permissions p ON rp.permission_id = p.id
WHERE rp.role_id = $1
  AND rp.permission_id = $2
  AND r.deleted_at IS NULL
  AND p.deleted_at IS NULL;

-- name: DeleteRolePermissionById :one
DELETE FROM role_permissions
WHERE id = $1
RETURNING *;

-- name: DeleteRolePermissionByRoleAndPermission :one
DELETE FROM role_permissions
WHERE role_id = $1 AND permission_id = $2
RETURNING *;

-- ============================================================================
-- BULK OPERATIONS
-- ============================================================================

-- name: BulkCreateRolePermissions :copyfrom
INSERT INTO role_permissions (id, role_id, permission_id)
VALUES ($1, $2, $3);

-- name: BulkAssignPermissionsToRole :many
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), $1, unnest($2::uuid[])
ON CONFLICT (role_id, permission_id) DO NOTHING
RETURNING *;

-- name: BulkAssignRolesToPermission :many
INSERT INTO role_permissions (id, role_id, permission_id)
SELECT gen_random_uuid(), unnest($1::uuid[]), $2
ON CONFLICT (role_id, permission_id) DO NOTHING
RETURNING *;

-- name: BulkDeleteRolePermissionByIds :exec
DELETE FROM role_permissions
WHERE id = ANY($1::uuid[]);

-- name: BulkRemovePermissionsFromRole :many
DELETE FROM role_permissions
WHERE role_id = $1 AND permission_id = ANY($2::uuid[])
RETURNING *;

-- name: BulkRemoveRolesFromPermission :many
DELETE FROM role_permissions
WHERE permission_id = $1 AND role_id = ANY($2::uuid[])
RETURNING *;

-- name: RemoveAllPermissionsFromRole :many
DELETE FROM role_permissions
WHERE role_id = $1
RETURNING *;

-- name: RemoveAllRolesFromPermission :many
DELETE FROM role_permissions
WHERE permission_id = $1
RETURNING *;

-- ============================================================================
-- LIST AND SEARCH OPERATIONS
-- ============================================================================

-- name: ListAllRolePermissions :many
SELECT rp.*,
       json_build_object(
           'id', r.id,
           'title_ru', r.title_ru,
           'title_en', r.title_en,
           'title_kk', r.title_kk,
           'value', r.value,
           'description_ru', r.description_ru,
           'description_en', r.description_en,
           'description_kk', r.description_kk,
           'created_at', r.created_at,
           'updated_at', r.updated_at,
           'deleted_at', r.deleted_at
       ) as role,
       json_build_object(
           'id', p.id,
           'title_ru', p.title_ru,
           'title_en', p.title_en,
           'title_kk', p.title_kk,
           'value', p.value,
           'description_ru', p.description_ru,
           'description_en', p.description_en,
           'description_kk', p.description_kk,
           'created_at', p.created_at,
           'updated_at', p.updated_at,
           'deleted_at', p.deleted_at
       ) as permission
FROM role_permissions rp
INNER JOIN roles r ON rp.role_id = r.id
INNER JOIN permissions p ON rp.permission_id = p.id
WHERE
    r.deleted_at IS NULL
    AND p.deleted_at IS NULL
    -- role_ids filter
    AND (
        sqlc.narg('role_ids')::uuid[] IS NULL OR
        rp.role_id = ANY(sqlc.narg('role_ids')::uuid[])
    )
    -- permission_ids filter
    AND (
        sqlc.narg('permission_ids')::uuid[] IS NULL OR
        rp.permission_id = ANY(sqlc.narg('permission_ids')::uuid[])
    )
    -- role_values filter
    AND (
        sqlc.narg('role_values')::text[] IS NULL OR
        r.value = ANY(sqlc.narg('role_values')::text[])
    )
    -- permission_values filter
    AND (
        sqlc.narg('permission_values')::text[] IS NULL OR
        p.value = ANY(sqlc.narg('permission_values')::text[])
    )
ORDER BY
    CASE WHEN sqlc.narg('sort_by') = 'created_at' AND sqlc.narg('sort_order') = 'ASC' THEN rp.created_at END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'created_at' AND sqlc.narg('sort_order') = 'DESC' THEN rp.created_at END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'role_value' AND sqlc.narg('sort_order') = 'ASC' THEN r.value END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'role_value' AND sqlc.narg('sort_order') = 'DESC' THEN r.value END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'permission_value' AND sqlc.narg('sort_order') = 'ASC' THEN p.value END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'permission_value' AND sqlc.narg('sort_order') = 'DESC' THEN p.value END DESC,
    rp.created_at DESC;

-- name: PaginateAllRolePermissions :many
SELECT rp.*,
       json_build_object(
           'id', r.id,
           'title_ru', r.title_ru,
           'title_en', r.title_en,
           'title_kk', r.title_kk,
           'value', r.value,
           'description_ru', r.description_ru,
           'description_en', r.description_en,
           'description_kk', r.description_kk,
           'created_at', r.created_at,
           'updated_at', r.updated_at,
           'deleted_at', r.deleted_at
       ) as role,
       json_build_object(
           'id', p.id,
           'title_ru', p.title_ru,
           'title_en', p.title_en,
           'title_kk', p.title_kk,
           'value', p.value,
           'description_ru', p.description_ru,
           'description_en', p.description_en,
           'description_kk', p.description_kk,
           'created_at', p.created_at,
           'updated_at', p.updated_at,
           'deleted_at', p.deleted_at
       ) as permission
FROM role_permissions rp
INNER JOIN roles r ON rp.role_id = r.id
INNER JOIN permissions p ON rp.permission_id = p.id
WHERE
    r.deleted_at IS NULL
    AND p.deleted_at IS NULL
    -- role_ids filter
    AND (
        sqlc.narg('role_ids')::uuid[] IS NULL OR
        rp.role_id = ANY(sqlc.narg('role_ids')::uuid[])
    )
    -- permission_ids filter
    AND (
        sqlc.narg('permission_ids')::uuid[] IS NULL OR
        rp.permission_id = ANY(sqlc.narg('permission_ids')::uuid[])
    )
    -- role_values filter
    AND (
        sqlc.narg('role_values')::text[] IS NULL OR
        r.value = ANY(sqlc.narg('role_values')::text[])
    )
    -- permission_values filter
    AND (
        sqlc.narg('permission_values')::text[] IS NULL OR
        p.value = ANY(sqlc.narg('permission_values')::text[])
    )
ORDER BY
    CASE WHEN sqlc.narg('sort_by') = 'created_at' AND sqlc.narg('sort_order') = 'ASC' THEN rp.created_at END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'created_at' AND sqlc.narg('sort_order') = 'DESC' THEN rp.created_at END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'role_value' AND sqlc.narg('sort_order') = 'ASC' THEN r.value END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'role_value' AND sqlc.narg('sort_order') = 'DESC' THEN r.value END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'permission_value' AND sqlc.narg('sort_order') = 'ASC' THEN p.value END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'permission_value' AND sqlc.narg('sort_order') = 'DESC' THEN p.value END DESC,
    rp.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountAllRolePermissions :one
SELECT COUNT(*)
FROM role_permissions rp
INNER JOIN roles r ON rp.role_id = r.id
INNER JOIN permissions p ON rp.permission_id = p.id
WHERE
    r.deleted_at IS NULL
    AND p.deleted_at IS NULL
    -- role_ids filter
    AND (
        sqlc.narg('role_ids')::uuid[] IS NULL OR
        rp.role_id = ANY(sqlc.narg('role_ids')::uuid[])
    )
    -- permission_ids filter
    AND (
        sqlc.narg('permission_ids')::uuid[] IS NULL OR
        rp.permission_id = ANY(sqlc.narg('permission_ids')::uuid[])
    )
    -- role_values filter
    AND (
        sqlc.narg('role_values')::text[] IS NULL OR
        r.value = ANY(sqlc.narg('role_values')::text[])
    )
    -- permission_values filter
    AND (
        sqlc.narg('permission_values')::text[] IS NULL OR
        p.value = ANY(sqlc.narg('permission_values')::text[])
    );

-- name: GetRolePermissions :many
SELECT p.*
FROM permissions p
INNER JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.role_id = $1 AND p.deleted_at IS NULL
ORDER BY p.created_at DESC;

-- name: GetPermissionRoles :many
SELECT r.*
FROM roles r
INNER JOIN role_permissions rp ON r.id = rp.role_id
WHERE rp.permission_id = $1 AND r.deleted_at IS NULL
ORDER BY r.created_at DESC;

-- ============================================================================
-- UTILITY OPERATIONS
-- ============================================================================

-- name: CheckRoleHasPermission :one
SELECT EXISTS(
    SELECT 1
    FROM role_permissions rp
    INNER JOIN roles r ON rp.role_id = r.id
    INNER JOIN permissions p ON rp.permission_id = p.id
    WHERE rp.role_id = $1
      AND rp.permission_id = $2
      AND r.deleted_at IS NULL
      AND p.deleted_at IS NULL
) as exists;

-- name: CheckRoleHasPermissionByValue :one
SELECT EXISTS(
    SELECT 1
    FROM role_permissions rp
    INNER JOIN roles r ON rp.role_id = r.id
    INNER JOIN permissions p ON rp.permission_id = p.id
    WHERE r.value = $1
      AND p.value = $2
      AND r.deleted_at IS NULL
      AND p.deleted_at IS NULL
) as exists;

-- name: CountRolePermissions :one
SELECT COUNT(*)
FROM role_permissions rp
INNER JOIN permissions p ON rp.permission_id = p.id
WHERE rp.role_id = $1
  AND p.deleted_at IS NULL;

-- name: CountPermissionRoles :one
SELECT COUNT(*)
FROM role_permissions rp
INNER JOIN roles r ON rp.role_id = r.id
WHERE rp.permission_id = $1
  AND r.deleted_at IS NULL;

-- ============================================================================
-- LEGACY QUERIES (For backward compatibility)
-- ============================================================================

-- name: AssignPermissionToRole :one
INSERT INTO role_permissions (id, role_id, permission_id)
VALUES ($1, $2, $3)
ON CONFLICT (role_id, permission_id) DO NOTHING
RETURNING *;

-- name: RemovePermissionFromRole :exec
DELETE FROM role_permissions
WHERE role_id = $1 AND permission_id = $2;

-- name: GetRolePermissionByID :one
SELECT * FROM role_permissions
WHERE id = $1;
