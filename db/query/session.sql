-- name: GetSessionsByCampaignAndCharacter :many
SELECT
    s.id, s.name,
    CASE
        WHEN sac.character_id IS NOT NULL THEN TRUE
        ELSE FALSE
        END AS isAllowed
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