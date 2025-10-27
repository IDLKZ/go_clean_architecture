-- ============================================================================
-- BASIC CRUD OPERATIONS
-- ============================================================================

-- name: CreateOnePermission :one
INSERT INTO permissions (id, title_ru, title_en, title_kk, description_ru, description_en, description_kk, value)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPermissionById :one
SELECT p.*,
       COALESCE(
           json_agg(
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
               )
           ) FILTER (WHERE r.id IS NOT NULL), '[]'
       ) as roles
FROM permissions p
LEFT JOIN role_permissions rp ON p.id = rp.permission_id
LEFT JOIN roles r ON rp.role_id = r.id AND r.deleted_at IS NULL
WHERE p.id = $1 AND p.deleted_at IS NULL
GROUP BY p.id;

-- name: GetPermissionByValue :one
SELECT p.*,
       COALESCE(
           json_agg(
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
               )
           ) FILTER (WHERE r.id IS NOT NULL), '[]'
       ) as roles
FROM permissions p
LEFT JOIN role_permissions rp ON p.id = rp.permission_id
LEFT JOIN roles r ON rp.role_id = r.id AND r.deleted_at IS NULL
WHERE p.value = $1 AND p.deleted_at IS NULL
GROUP BY p.id;

-- name: UpdatePermissionById :one
UPDATE permissions
SET title_ru = $2,
    title_en = $3,
    title_kk = $4,
    description_ru = $5,
    description_en = $6,
    description_kk = $7,
    value = $8,
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeletePermissionById :one
UPDATE permissions
SET deleted_at = now(),
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: HardDeletePermissionById :exec
DELETE FROM permissions
WHERE id = $1;

-- ============================================================================
-- BULK OPERATIONS
-- ============================================================================

-- name: BulkCreatePermissions :copyfrom
INSERT INTO permissions (id, title_ru, title_en, title_kk, description_ru, description_en, description_kk, value)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: BulkUpdatePermissions :batchmany
UPDATE permissions
SET title_ru = $2,
    title_en = $3,
    title_kk = $4,
    description_ru = $5,
    description_en = $6,
    description_kk = $7,
    value = $8,
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: BulkDeletePermissionByIds :many
UPDATE permissions
SET deleted_at = now(),
    updated_at = now()
WHERE id = ANY($1::uuid[]) AND deleted_at IS NULL
RETURNING *;

-- name: BulkHardDeletePermissionByIds :exec
DELETE FROM permissions
WHERE id = ANY($1::uuid[]);

-- ============================================================================
-- LIST AND SEARCH OPERATIONS
-- ============================================================================

-- name: ListAllPermissions :many
SELECT p.*,
       COALESCE(
           json_agg(
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
               )
           ) FILTER (WHERE r.id IS NOT NULL), '[]'
       ) as roles
FROM permissions p
LEFT JOIN role_permissions rp ON p.id = rp.permission_id
LEFT JOIN roles r ON rp.role_id = r.id AND r.deleted_at IS NULL
WHERE
    -- show_deleted filter
    (CASE WHEN sqlc.narg('show_deleted')::boolean THEN TRUE ELSE p.deleted_at IS NULL END)
    -- search filter (title_ru, title_en, title_kk, description_ru, description_en, description_kk, value)
    AND (
        sqlc.narg('search')::text IS NULL OR
        p.title_ru ILIKE '%' || sqlc.narg('search') || '%' OR
        p.title_en ILIKE '%' || sqlc.narg('search') || '%' OR
        p.title_kk ILIKE '%' || sqlc.narg('search') || '%' OR
        p.description_ru ILIKE '%' || sqlc.narg('search') || '%' OR
        p.description_en ILIKE '%' || sqlc.narg('search') || '%' OR
        p.description_kk ILIKE '%' || sqlc.narg('search') || '%' OR
        p.value ILIKE '%' || sqlc.narg('search') || '%'
    )
    -- values filter
    AND (
        sqlc.narg('values')::text[] IS NULL OR
        p.value = ANY(sqlc.narg('values')::text[])
    )
    -- ids filter
    AND (
        sqlc.narg('ids')::uuid[] IS NULL OR
        p.id = ANY(sqlc.narg('ids')::uuid[])
    )
GROUP BY p.id
ORDER BY
    CASE WHEN sqlc.narg('sort_by') = 'created_at' AND sqlc.narg('sort_order') = 'ASC' THEN p.created_at END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'created_at' AND sqlc.narg('sort_order') = 'DESC' THEN p.created_at END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'updated_at' AND sqlc.narg('sort_order') = 'ASC' THEN p.updated_at END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'updated_at' AND sqlc.narg('sort_order') = 'DESC' THEN p.updated_at END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'title_ru' AND sqlc.narg('sort_order') = 'ASC' THEN p.title_ru END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'title_ru' AND sqlc.narg('sort_order') = 'DESC' THEN p.title_ru END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'value' AND sqlc.narg('sort_order') = 'ASC' THEN p.value END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'value' AND sqlc.narg('sort_order') = 'DESC' THEN p.value END DESC,
    p.created_at DESC;

-- name: PaginateAllPermissions :many
SELECT p.*,
       COALESCE(
           json_agg(
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
               )
           ) FILTER (WHERE r.id IS NOT NULL), '[]'
       ) as roles
FROM permissions p
LEFT JOIN role_permissions rp ON p.id = rp.permission_id
LEFT JOIN roles r ON rp.role_id = r.id AND r.deleted_at IS NULL
WHERE
    -- show_deleted filter
    (CASE WHEN sqlc.narg('show_deleted')::boolean THEN TRUE ELSE p.deleted_at IS NULL END)
    -- search filter (title_ru, title_en, title_kk, description_ru, description_en, description_kk, value)
    AND (
        sqlc.narg('search')::text IS NULL OR
        p.title_ru ILIKE '%' || sqlc.narg('search') || '%' OR
        p.title_en ILIKE '%' || sqlc.narg('search') || '%' OR
        p.title_kk ILIKE '%' || sqlc.narg('search') || '%' OR
        p.description_ru ILIKE '%' || sqlc.narg('search') || '%' OR
        p.description_en ILIKE '%' || sqlc.narg('search') || '%' OR
        p.description_kk ILIKE '%' || sqlc.narg('search') || '%' OR
        p.value ILIKE '%' || sqlc.narg('search') || '%'
    )
    -- values filter
    AND (
        sqlc.narg('values')::text[] IS NULL OR
        p.value = ANY(sqlc.narg('values')::text[])
    )
    -- ids filter
    AND (
        sqlc.narg('ids')::uuid[] IS NULL OR
        p.id = ANY(sqlc.narg('ids')::uuid[])
    )
GROUP BY p.id
ORDER BY
    CASE WHEN sqlc.narg('sort_by') = 'created_at' AND sqlc.narg('sort_order') = 'ASC' THEN p.created_at END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'created_at' AND sqlc.narg('sort_order') = 'DESC' THEN p.created_at END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'updated_at' AND sqlc.narg('sort_order') = 'ASC' THEN p.updated_at END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'updated_at' AND sqlc.narg('sort_order') = 'DESC' THEN p.updated_at END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'title_ru' AND sqlc.narg('sort_order') = 'ASC' THEN p.title_ru END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'title_ru' AND sqlc.narg('sort_order') = 'DESC' THEN p.title_ru END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'value' AND sqlc.narg('sort_order') = 'ASC' THEN p.value END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'value' AND sqlc.narg('sort_order') = 'DESC' THEN p.value END DESC,
    p.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountAllPermissions :one
SELECT COUNT(DISTINCT p.id)
FROM permissions p
WHERE
    -- show_deleted filter
    (CASE WHEN sqlc.narg('show_deleted')::boolean THEN TRUE ELSE p.deleted_at IS NULL END)
    -- search filter
    AND (
        sqlc.narg('search')::text IS NULL OR
        p.title_ru ILIKE '%' || sqlc.narg('search') || '%' OR
        p.title_en ILIKE '%' || sqlc.narg('search') || '%' OR
        p.title_kk ILIKE '%' || sqlc.narg('search') || '%' OR
        p.description_ru ILIKE '%' || sqlc.narg('search') || '%' OR
        p.description_en ILIKE '%' || sqlc.narg('search') || '%' OR
        p.description_kk ILIKE '%' || sqlc.narg('search') || '%' OR
        p.value ILIKE '%' || sqlc.narg('search') || '%'
    )
    -- values filter
    AND (
        sqlc.narg('values')::text[] IS NULL OR
        p.value = ANY(sqlc.narg('values')::text[])
    )
    -- ids filter
    AND (
        sqlc.narg('ids')::uuid[] IS NULL OR
        p.id = ANY(sqlc.narg('ids')::uuid[])
    );

-- ============================================================================
-- LEGACY QUERIES (For backward compatibility)
-- ============================================================================

-- name: GetPermissionWithRoles :one
SELECT p.*,
       COALESCE(
           json_agg(
               json_build_object(
                   'id', r.id,
                   'title_ru', r.title_ru,
                   'title_en', r.title_en,
                   'title_kk', r.title_kk,
                   'value', r.value,
                   'description_ru', r.description_ru,
                   'description_en', r.description_en,
                   'description_kk', r.description_kk
               )
           ) FILTER (WHERE r.id IS NOT NULL), '[]'
       ) as roles
FROM permissions p
LEFT JOIN role_permissions rp ON p.id = rp.permission_id
LEFT JOIN roles r ON rp.role_id = r.id AND r.deleted_at IS NULL
WHERE p.id = $1 AND p.deleted_at IS NULL
GROUP BY p.id;
