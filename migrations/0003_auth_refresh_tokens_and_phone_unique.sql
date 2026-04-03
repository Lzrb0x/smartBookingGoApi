BEGIN;

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_id VARCHAR(100) NOT NULL UNIQUE,
    token_hash TEXT NOT NULL,
    expires_on TIMESTAMPTZ NOT NULL,
    created_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_on TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS refresh_tokens_user_id_idx ON refresh_tokens (user_id);
CREATE INDEX IF NOT EXISTS refresh_tokens_token_id_active_idx ON refresh_tokens (token_id) WHERE revoked_on IS NULL;

DROP INDEX IF EXISTS users_phone_unique_idx;
CREATE UNIQUE INDEX users_phone_unique_idx ON users (phone) WHERE phone IS NOT NULL AND phone != '';

COMMIT;
