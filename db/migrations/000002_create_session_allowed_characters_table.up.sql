CREATE TABLE "session_allowed_characters" (
  session_id BIGINT NOT NULL REFERENCES "sessions"(id) ON DELETE CASCADE,
  character_id BIGINT NOT NULL,
  PRIMARY KEY (session_id, character_id)
);
