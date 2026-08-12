-- +goose Up
CREATE TABLE users (
    id                    bigserial    PRIMARY KEY,
    created_at            timestamptz,
    updated_at            timestamptz,
    deleted_at            timestamptz,

    username              varchar(50)  NOT NULL,
    email                 varchar(100) NOT NULL,
    name                  varchar(255) NOT NULL,
    password              varchar(60)  NOT NULL,
    phone                 varchar(20)  NOT NULL,
    birth_date            date         NOT NULL,

    is_online             boolean      NOT NULL,
    last_login_at         timestamptz,
    last_logout_at        timestamptz,
    failed_login_attempts bigint       NOT NULL,
    locked_until          timestamptz,
    is_account_disabled   boolean      NOT NULL
);

-- Nomes exatos esperados pelo switch em internal/repository/user_repository.go
CREATE UNIQUE INDEX idx_users_username   ON users (username);
CREATE UNIQUE INDEX idx_users_email      ON users (email);
CREATE        INDEX idx_users_deleted_at ON users (deleted_at);

-- +goose Down
DROP TABLE users;
