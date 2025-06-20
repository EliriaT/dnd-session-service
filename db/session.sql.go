// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: session.sql

package db

import (
	"context"
	"time"

	"github.com/lib/pq"
)

const addSessionAllowedCharacter = `-- name: AddSessionAllowedCharacter :exec
INSERT INTO "session_allowed_characters" (session_id, character_id)
SELECT $1, unnest($2::bigint[])
ON CONFLICT DO NOTHING
`

type AddSessionAllowedCharacterParams struct {
	SessionID int64   `json:"sessionId"`
	Column2   []int64 `json:"column2"`
}

func (q *Queries) AddSessionAllowedCharacter(ctx context.Context, arg AddSessionAllowedCharacterParams) error {
	_, err := q.db.ExecContext(ctx, addSessionAllowedCharacter, arg.SessionID, pq.Array(arg.Column2))
	return err
}

const createSession = `-- name: CreateSession :one
INSERT INTO "sessions" (name, campaign_id, map_id)
VALUES ($1, $2, $3)
RETURNING id, name, campaign_id, map_id, is_active
`

type CreateSessionParams struct {
	Name       string `json:"name"`
	CampaignID int64  `json:"campaignId"`
	MapID      int64  `json:"mapId"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, createSession, arg.Name, arg.CampaignID, arg.MapID)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CampaignID,
		&i.MapID,
		&i.IsActive,
	)
	return i, err
}

const getCharactersBySession = `-- name: GetCharactersBySession :many
SELECT session_id, char_id, x_pos, y_pos, is_visible, modification_date
FROM session_characters_position
WHERE session_id = $1
`

func (q *Queries) GetCharactersBySession(ctx context.Context, sessionID int64) ([]SessionCharactersPosition, error) {
	rows, err := q.db.QueryContext(ctx, getCharactersBySession, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SessionCharactersPosition
	for rows.Next() {
		var i SessionCharactersPosition
		if err := rows.Scan(
			&i.SessionID,
			&i.CharID,
			&i.XPos,
			&i.YPos,
			&i.IsVisible,
			&i.ModificationDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getObjectsBySession = `-- name: GetObjectsBySession :many
SELECT session_id, object_id, x_pos, y_pos, is_visible, modification_date
FROM session_objects_position
WHERE session_id = $1
`

func (q *Queries) GetObjectsBySession(ctx context.Context, sessionID int64) ([]SessionObjectsPosition, error) {
	rows, err := q.db.QueryContext(ctx, getObjectsBySession, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SessionObjectsPosition
	for rows.Next() {
		var i SessionObjectsPosition
		if err := rows.Scan(
			&i.SessionID,
			&i.ObjectID,
			&i.XPos,
			&i.YPos,
			&i.IsVisible,
			&i.ModificationDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSessionByID = `-- name: GetSessionByID :one
SELECT id, name, campaign_id, map_id, is_active
FROM sessions
WHERE id = $1
`

func (q *Queries) GetSessionByID(ctx context.Context, id int64) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSessionByID, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CampaignID,
		&i.MapID,
		&i.IsActive,
	)
	return i, err
}

const getSessionsByCampaignAndCharacter = `-- name: GetSessionsByCampaignAndCharacter :many
SELECT
    s.id, s.name, s.is_active,
    CASE
        WHEN sac.character_id IS NOT NULL THEN TRUE
        ELSE FALSE
        END AS is_allowed
FROM "sessions" s
         LEFT JOIN "session_allowed_characters" sac
                   ON s.id = sac.session_id AND sac.character_id = $2
WHERE s.campaign_id = $1
`

type GetSessionsByCampaignAndCharacterParams struct {
	CampaignID  int64 `json:"campaignId"`
	CharacterID int64 `json:"characterId"`
}

type GetSessionsByCampaignAndCharacterRow struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	IsActive  bool   `json:"isActive"`
	IsAllowed bool   `json:"isAllowed"`
}

func (q *Queries) GetSessionsByCampaignAndCharacter(ctx context.Context, arg GetSessionsByCampaignAndCharacterParams) ([]GetSessionsByCampaignAndCharacterRow, error) {
	rows, err := q.db.QueryContext(ctx, getSessionsByCampaignAndCharacter, arg.CampaignID, arg.CharacterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSessionsByCampaignAndCharacterRow
	for rows.Next() {
		var i GetSessionsByCampaignAndCharacterRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.IsActive,
			&i.IsAllowed,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setSessionActive = `-- name: SetSessionActive :exec
UPDATE sessions
SET is_active = TRUE
WHERE id = $1
`

func (q *Queries) SetSessionActive(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, setSessionActive, id)
	return err
}

const upsertCharacterPosition = `-- name: UpsertCharacterPosition :exec
INSERT INTO session_characters_position (
    session_id, char_id, x_pos, y_pos, is_visible, modification_date
)
VALUES ($1, $2, $3, $4, $5, $6)
    ON CONFLICT (session_id, char_id) DO UPDATE
    SET x_pos = EXCLUDED.x_pos,
    y_pos = EXCLUDED.y_pos,
    is_visible = EXCLUDED.is_visible,
    modification_date = EXCLUDED.modification_date
`

type UpsertCharacterPositionParams struct {
	SessionID        int64     `json:"sessionId"`
	CharID           int64     `json:"charId"`
	XPos             int32     `json:"xPos"`
	YPos             int32     `json:"yPos"`
	IsVisible        bool      `json:"isVisible"`
	ModificationDate time.Time `json:"modificationDate"`
}

func (q *Queries) UpsertCharacterPosition(ctx context.Context, arg UpsertCharacterPositionParams) error {
	_, err := q.db.ExecContext(ctx, upsertCharacterPosition,
		arg.SessionID,
		arg.CharID,
		arg.XPos,
		arg.YPos,
		arg.IsVisible,
		arg.ModificationDate,
	)
	return err
}

const upsertObjectPosition = `-- name: UpsertObjectPosition :exec
INSERT INTO session_objects_position (
    session_id, object_id, x_pos, y_pos, is_visible, modification_date
)
VALUES ($1, $2, $3, $4, $5, $6)
    ON CONFLICT (session_id, object_id) DO UPDATE
       SET x_pos = EXCLUDED.x_pos,
       y_pos = EXCLUDED.y_pos,
       is_visible = EXCLUDED.is_visible,
       modification_date = EXCLUDED.modification_date
`

type UpsertObjectPositionParams struct {
	SessionID        int64     `json:"sessionId"`
	ObjectID         int64     `json:"objectId"`
	XPos             int32     `json:"xPos"`
	YPos             int32     `json:"yPos"`
	IsVisible        bool      `json:"isVisible"`
	ModificationDate time.Time `json:"modificationDate"`
}

func (q *Queries) UpsertObjectPosition(ctx context.Context, arg UpsertObjectPositionParams) error {
	_, err := q.db.ExecContext(ctx, upsertObjectPosition,
		arg.SessionID,
		arg.ObjectID,
		arg.XPos,
		arg.YPos,
		arg.IsVisible,
		arg.ModificationDate,
	)
	return err
}
