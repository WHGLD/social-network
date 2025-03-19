-- +migrate Up

BEGIN;

CREATE TYPE gender AS ENUM ('male', 'female');

CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(100),
    second_name VARCHAR(100),
    birthday VARCHAR(100),
    city VARCHAR(100),
    sex gender,
    biography VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE users IS 'Таблица для хранения данных пользователей';
COMMENT ON COLUMN users.password_hash IS 'Хэш пароля пользователя';

COMMIT;