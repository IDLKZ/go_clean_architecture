-- ============================================================================
-- BASIC CRUD OPERATIONS
-- ============================================================================

-- name: CreateOneRole :one
INSERT INTO roles (id, title_ru, title_en, title_kk, description_ru, description_en, description_kk, value)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetRoleById :one
SELECT r.*,
       COALESCE(
           json_agg(
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
               )
           ) FILTER (WHERE p.id IS NOT NULL), '[]'
       ) as permissions
FROM roles r
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id AND p.deleted_at IS NULL
WHERE r.id = $1 AND r.deleted_at IS NULL
GROUP BY r.id;

-- name: GetRoleByValue :one
SELECT r.*,
       COALESCE(
           json_agg(
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
               )
           ) FILTER (WHERE p.id IS NOT NULL), '[]'
       ) as permissions
FROM roles r
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id AND p.deleted_at IS NULL
WHERE r.value = $1 AND r.deleted_at IS NULL
GROUP BY r.id;

-- name: UpdateRoleById :one
UPDATE roles
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

-- name: DeleteRoleById :one
UPDATE roles
SET deleted_at = now(),
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: HardDeleteRoleById :exec
DELETE FROM roles
WHERE id = $1;

-- ============================================================================
-- BULK OPERATIONS
-- ============================================================================

-- name: BulkCreateRoles :copyfrom
INSERT INTO roles (id, title_ru, title_en, title_kk, description_ru, description_en, description_kk, value)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: BulkUpdateRoles :batchmany
UPDATE roles
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

-- name: BulkDeleteRoleByIds :many
UPDATE roles
SET deleted_at = now(),
    updated_at = now()
WHERE id = ANY($1::uuid[]) AND deleted_at IS NULL
RETURNING *;

-- name: BulkHardDeleteRoleByIds :exec
DELETE FROM roles
WHERE id = ANY($1::uuid[]);

-- ============================================================================
-- LIST AND SEARCH OPERATIONS
-- ============================================================================

-- name: ListAllRoles :many
SELECT r.*,
       COALESCE(
           json_agg(
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
               )
           ) FILTER (WHERE p.id IS NOT NULL), '[]'
       ) as permissions
FROM roles r
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id AND p.deleted_at IS NULL
WHERE
    -- show_deleted filter
    (CASE WHEN sqlc.narg('show_deleted')::boolean THEN TRUE ELSE r.deleted_at IS NULL END)
    -- search filter (title_ru, title_en, title_kk, description_ru, description_en, description_kk, value)
    AND (
        sqlc.narg('search')::text IS NULL OR
        r.title_ru ILIKE '%' || sqlc.narg('search') || '%' OR
        r.title_en ILIKE '%' || sqlc.narg('search') || '%' OR
        r.title_kk ILIKE '%' || sqlc.narg('search') || '%' OR
        r.description_ru ILIKE '%' || sqlc.narg('search') || '%' OR
        r.description_en ILIKE '%' || sqlc.narg('search') || '%' OR
        r.description_kk ILIKE '%' || sqlc.narg('search') || '%' OR
        r.value ILIKE '%' || sqlc.narg('search') || '%'
    )
    -- values filter
    AND (
        sqlc.narg('values')::text[] IS NULL OR
        r.value = ANY(sqlc.narg('values')::text[])
    )
    -- ids filter
    AND (
        sqlc.narg('ids')::uuid[] IS NULL OR
        r.id = ANY(sqlc.narg('ids')::uuid[])
    )
GROUP BY r.id
ORDER BY
    CASE WHEN sqlc.narg('sort_by') = 'created_at' AND sqlc.narg('sort_order') = 'ASC' THEN r.created_at END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'created_at' AND sqlc.narg('sort_order') = 'DESC' THEN r.created_at END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'updated_at' AND sqlc.narg('sort_order') = 'ASC' THEN r.updated_at END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'updated_at' AND sqlc.narg('sort_order') = 'DESC' THEN r.updated_at END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'title_ru' AND sqlc.narg('sort_order') = 'ASC' THEN r.title_ru END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'title_ru' AND sqlc.narg('sort_order') = 'DESC' THEN r.title_ru END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'value' AND sqlc.narg('sort_order') = 'ASC' THEN r.value END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'value' AND sqlc.narg('sort_order') = 'DESC' THEN r.value END DESC,
    r.created_at DESC;

-- name: PaginateAllRoles :many
SELECT r.*,
       COALESCE(
           json_agg(
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
               )
           ) FILTER (WHERE p.id IS NOT NULL), '[]'
       ) as permissions
FROM roles r
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id AND p.deleted_at IS NULL
WHERE
    -- show_deleted filter
    (CASE WHEN sqlc.narg('show_deleted')::boolean THEN TRUE ELSE r.deleted_at IS NULL END)
    -- search filter (title_ru, title_en, title_kk, description_ru, description_en, description_kk, value)
    AND (
        sqlc.narg('search')::text IS NULL OR
        r.title_ru ILIKE '%' || sqlc.narg('search') || '%' OR
        r.title_en ILIKE '%' || sqlc.narg('search') || '%' OR
        r.title_kk ILIKE '%' || sqlc.narg('search') || '%' OR
        r.description_ru ILIKE '%' || sqlc.narg('search') || '%' OR
        r.description_en ILIKE '%' || sqlc.narg('search') || '%' OR
        r.description_kk ILIKE '%' || sqlc.narg('search') || '%' OR
        r.value ILIKE '%' || sqlc.narg('search') || '%'
    )
    -- values filter
    AND (
        sqlc.narg('values')::text[] IS NULL OR
        r.value = ANY(sqlc.narg('values')::text[])
    )
    -- ids filter
    AND (
        sqlc.narg('ids')::uuid[] IS NULL OR
        r.id = ANY(sqlc.narg('ids')::uuid[])
    )
GROUP BY r.id
ORDER BY
    CASE WHEN sqlc.narg('sort_by') = 'created_at' AND sqlc.narg('sort_order') = 'ASC' THEN r.created_at END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'created_at' AND sqlc.narg('sort_order') = 'DESC' THEN r.created_at END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'updated_at' AND sqlc.narg('sort_order') = 'ASC' THEN r.updated_at END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'updated_at' AND sqlc.narg('sort_order') = 'DESC' THEN r.updated_at END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'title_ru' AND sqlc.narg('sort_order') = 'ASC' THEN r.title_ru END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'title_ru' AND sqlc.narg('sort_order') = 'DESC' THEN r.title_ru END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'value' AND sqlc.narg('sort_order') = 'ASC' THEN r.value END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'value' AND sqlc.narg('sort_order') = 'DESC' THEN r.value END DESC,
    r.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountAllRoles :one
SELECT COUNT(DISTINCT r.id)
FROM roles r
WHERE
    -- show_deleted filter
    (CASE WHEN sqlc.narg('show_deleted')::boolean THEN TRUE ELSE r.deleted_at IS NULL END)
    -- search filter
    AND (
        sqlc.narg('search')::text IS NULL OR
        r.title_ru ILIKE '%' || sqlc.narg('search') || '%' OR
        r.title_en ILIKE '%' || sqlc.narg('search') || '%' OR
        r.title_kk ILIKE '%' || sqlc.narg('search') || '%' OR
        r.description_ru ILIKE '%' || sqlc.narg('search') || '%' OR
        r.description_en ILIKE '%' || sqlc.narg('search') || '%' OR
        r.description_kk ILIKE '%' || sqlc.narg('search') || '%' OR
        r.value ILIKE '%' || sqlc.narg('search') || '%'
    )
    -- values filter
    AND (
        sqlc.narg('values')::text[] IS NULL OR
        r.value = ANY(sqlc.narg('values')::text[])
    )
    -- ids filter
    AND (
        sqlc.narg('ids')::uuid[] IS NULL OR
        r.id = ANY(sqlc.narg('ids')::uuid[])
    );

-- ============================================================================
-- LEGACY QUERIES (For backward compatibility)
-- ============================================================================

-- name: GetRoleWithPermissions :one
SELECT r.*,
       COALESCE(
           json_agg(
               json_build_object(
                   'id', p.id,
                   'title_ru', p.title_ru,
                   'title_en', p.title_en,
                   'title_kk', p.title_kk,
                   'value', p.value,
                   'description_ru', p.description_ru,
                   'description_en', p.description_en,
                   'description_kk', p.description_kk
               )
           ) FILTER (WHERE p.id IS NOT NULL), '[]'
       ) as permissions
FROM roles r
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id AND p.deleted_at IS NULL
WHERE r.id = $1 AND r.deleted_at IS NULL
GROUP BY r.id;
