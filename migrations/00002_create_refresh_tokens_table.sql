-- +goose Up
CREATE TABLE refresh_tokens (
    id         bigserial   PRIMARY KEY,
    created_at timestamptz NOT NULL,
    user_id    bigint      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash bytea       NOT NULL,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz
);

CREATE UNIQUE INDEX idx_refresh_tokens_token_hash ON refresh_tokens (token_hash);
CREATE        INDEX idx_refresh_tokens_user_id    ON refresh_tokens (user_id);

-- +goose Down
DROP TABLE refresh_tokens;
