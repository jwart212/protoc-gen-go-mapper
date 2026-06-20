
-- name: CreateOutboxEvent :one
INSERT INTO schm_pos.outbox (aggregate_id, topic, key, payload, headers, attempts, created_at) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING *;

-- name: ChangeStatusEvent :one
UPDATE schm_pos.outbox SET status=$1 WHERE id = $2 RETURNING *;

-- name: ChangeAttemptCount :one
UPDATE schm_pos.outbox SET attempts=$1 WHERE id=$2 RETURNING *;

-- name: ChangeRetryCount :one
UPDATE schm_pos.outbox SET retry_count=$1 WHERE id=$2 RETURNING *;

-- name: UpdateOutbox :one
UPDATE schm_pos.outbox
SET 
    status          = COALESCE(sqlc.narg(status), status),
    attempts        = COALESCE(sqlc.narg(attempts), attempts),
    retry_count     = COALESCE(sqlc.narg(retry_count), retry_count),
    error_msg       = COALESCE(sqlc.narg(error_msg), error_msg),
    retry_at        = COALESCE(sqlc.narg(retry_at), retry_at),
    next_attempt_at = COALESCE(sqlc.narg(next_attempt_at), next_attempt_at),
    updated_at      = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;


-- name: OutboxLists :many
SELECT 
    id,
    aggregate_id,
    topic,
    key,
    payload,
    headers,
    attempts,
    status,
    error_msg,
    retry_count,
    retry_at,
    created_at,
    updated_at,
    next_attempt_at
FROM schm_pos.outbox;

-- name: GetOutbox :one
SELECT 
    id,
    aggregate_id,
    topic,
    key,
    payload,
    headers,
    attempts,
    status,
    error_msg,
    retry_count,
    retry_at,
    created_at,
    updated_at,
    next_attempt_at
FROM schm_pos.outbox
WHERE
    id = $1;

-- name: GetPendingOutbox :many
SELECT
    id,
    aggregate_id,
    topic,
    key,
    payload,
    headers,
    attempts,
    status,
    error_msg,
    retry_count,
    retry_at,
    created_at,
    updated_at,
    next_attempt_at
FROM schm_pos.outbox
WHERE
    status = 'PENDING'
    AND (
        next_attempt_at IS NULL
        OR next_attempt_at <= NOW()
    )
ORDER BY created_at ASC
LIMIT $1;

-- name: DestroyOutbox :one
DELETE FROM schm_pos.outbox WHERE id=$1 RETURNING *;


-- name: CreateUOM :one
INSERT INTO schm_pos.uoms (code, name, symbol, description, created_at) VALUES($1,$2,$3,$4,$5) RETURNING *;

-- name: UpdatedUOM :one
UPDATE schm_pos.uoms
SET
    code        = COALESCE(sqlc.narg(code), code),
    name        = COALESCE(sqlc.narg(name), name),
    symbol      = COALESCE(sqlc.narg(symbol), symbol),
    description = COALESCE(sqlc.narg(description), description),
    updated_at  = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: ListsUOMS :many
SELECT 
    code,
    name,
    symbol,
    description,
    created_at,
    updated_at
FROM schm_pos.uoms WHERE deleted_at IS NULL;

-- name: GetUOM :one
SELECT 
    code,
    name,
    symbol,
    description,
    created_at,
    updated_at
FROM schm_pos.uoms
WHERE id = sqlc.arg(id) AND deleted_at IS NULL;

-- name: DeleteUOM :one
UPDATE schm_pos.uoms
SET
    deleted_at = NOW(),
    updated_at = NOW()
WHERE id=$1 RETURNING *;

-- name: RestoreUOM :one
UPDATE schm_pos.uoms
SET
    deleted_at = NULL,
    updated_at = NOW()
WHERE id=$1 RETURNING *;

-- name: ExistsUOMByID :one
SELECT EXISTS (
    SELECT 1 
        FROM schm_pos.uoms
    WHERE 
        id = sqlc.arg(id)
    AND 
        deleted_at IS NULL
);

-- name: ExistsUOMByCode :one
SELECT EXISTS (
    SELECT 1
        FROM schm_pos.uoms
    WHERE
        code = sqlc.arg(code)
    AND 
        deleted_at IS NULL
);

-- name: ExistsUOMBySymbol :one
SELECT EXISTS (
    SELECT 1
        FROM schm_pos.uoms
    WHERE
        symbol = sqlc.arg(symbol)
    AND 
        deleted_at IS NULL
);

-- name: ExistsUOMCodeExcludeID :one
SELECT EXISTS (
    SELECT 1
    FROM schm_pos.uoms
    WHERE
        code = sqlc.arg(code)
        AND id <> sqlc.arg(id)
        AND deleted_at IS NULL
);

-- name: BatchDeleteUOM :many
UPDATE schm_pos.uoms
SET
    deleted_at = NOW(),
    updated_at = NOW()
WHERE
    id = ANY(sqlc.arg(ids)::uuid[])
RETURNING *;

-- name: BatchRestoreUOM :many
UPDATE schm_pos.uoms
SET
    deleted_at = NULL,
    updated_at = NOW()
WHERE
    id = ANY(sqlc.arg(ids)::uuid[])
RETURNING *;


-- name: CreateItemCategories :one
INSERT INTO schm_pos.item_categories (tenant_id,parent_id,code,name,description,created_at)
VALUES ($1,$2,$3,$4,$5,NOW()) RETURNING *;

-- name: UpdateItemCategories :one
UPDATE schm_pos.item_categories
SET
    tenant_id   = COALESCE(sqlc.narg(tenant_id), tenant_id),
    parent_id   = COALESCE(sqlc.narg(parent_id), parent_id),
    code        = COALESCE(sqlc.narg(code), code),
    name        = COALESCE(sqlc.narg(name), name),
    description = COALESCE(sqlc.narg(description), description),
    updated_at  = NOW()
WHERE id = sqlc.arg(id) RETURNING *;

-- name: DeleteItemCategories :one
UPDATE schm_pos.item_categories
SET
    updated_at  = NOW(),
    deleted_at  = NOW(),
    deleted_by  = $1
WHERE id = sqlc.arg(id) RETURNING *;

-- name: RestoreItemCategories :one
UPDATE schm_pos.item_categories
SET
    updated_at  = NOW(),
    deleted_at  = NULL
WHERE id = sqlc.arg(id) RETURNING *;

-- name: ListsItemCategories :many
SELECT 
    id,
    tenant_id,
    parent_id,
    code,
    name,
    description,
    created_at,
    updated_at
FROM schm_pos.item_categories 
WHERE 
    deleted_at IS NULL
    AND (
        sqlc.narg('tenant_id')::uuid IS NULL
        OR tenant_id = sqlc.narg('tenant_id')::uuid
    )
ORDER BY name;


-- name: ListsItemCategoriesTree :many
WITH RECURSIVE category_tree AS (
    SELECT
        ic.id,
        ic.tenant_id,
        ic.parent_id,
        ic.code,
        ic.name,
        ic.description,
        ic.created_at,
        ic.updated_at,
        ic.deleted_at,
        0::int AS level,
        ARRAY[id::text] AS path
    FROM schm_pos.item_categories as ic
    WHERE
        ic.parent_id IS NULL
        AND ic.tenant_id = $1
        AND ic.deleted_at IS NULL

    UNION ALL

    SELECT
        c.id,
        c.parent_id,
        c.code,
        c.name,
        c.description,
        c.created_at,
        c.updated_at,
        c.deleted_at,
        ct.level + 1,
        ct.path || c.id::text
    FROM schm_pos.item_categories c
    INNER JOIN category_tree ct
        ON c.parent_id = ct.id
    WHERE
        c.deleted_at IS NULL
)
SELECT *
FROM category_tree
ORDER BY path;

-- name: GetItemCategories :one
SELECT 
    id,
    tenant_id,
    parent_id,
    code,
    name,
    description,
    created_at,
    updated_at
FROM schm_pos.item_categories
WHERE id=sqlc.arg('id') AND deleted_at IS NULL;