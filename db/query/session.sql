-- name: GetSessionsByCampaignAndCharacter :many
SELECT
    s.id, s.name,
    CASE
        WHEN sac.character_id IS NOT NULL THEN TRUE
        ELSE FALSE
        END AS is_allowed
FROM "sessions" s
         LEFT JOIN "session_allowed_characters" sac
                   ON s.id = sac.session_id AND sac.character_id = $2
WHERE s.campaign_id = $1;

-- name: CreateSession :one
INSERT INTO "sessions" (name, campaign_id, map_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: AddSessionAllowedCharacter :exec
INSERT INTO "session_allowed_characters" (session_id, character_id)
SELECT $1, unnest($2::bigint[])
ON CONFLICT DO NOTHING;

-- name: GetObjectsBySession :many
SELECT *
FROM session_objects_position
WHERE session_id = $1;

SELECT *
FROM session_characters_position
WHERE session_id = $1;

-- name: UpsertCharacterPosition :exec
INSERT INTO session_characters_position (
    session_id, char_id, x_pos, y_pos, is_visible, modification_date
)
VALUES ($1, $2, $3, $4, $5, $6)
    ON CONFLICT (session_id, char_id) DO UPDATE
    SET x_pos = EXCLUDED.x_pos,
    y_pos = EXCLUDED.y_pos,
    is_visible = EXCLUDED.is_visible,
    modification_date = EXCLUDED.modification_date;

-- name: UpsertObjectPosition :exec
INSERT INTO session_objects_position (
    session_id, object_id, x_pos, y_pos, is_visible, modification_date
)
VALUES ($1, $2, $3, $4, $5, $6)
    ON CONFLICT (session_id, object_id) DO UPDATE
       SET x_pos = EXCLUDED.x_pos,
       y_pos = EXCLUDED.y_pos,
       is_visible = EXCLUDED.is_visible,
       modification_date = EXCLUDED.modification_date;