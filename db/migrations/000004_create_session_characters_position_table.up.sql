CREATE TABLE "session_characters_position" (
  session_id BIGINT NOT NULL REFERENCES "sessions"(id) ON DELETE CASCADE,
  char_id BIGINT NOT NULL,
  x_pos INTEGER NOT NULL,
  y_pos INTEGER NOT NULL,
  is_visible BOOLEAN NOT NULL DEFAULT TRUE,
  modification_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (session_id, char_id)
);