-- +migrate Up

BEGIN;

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS users_names_ind ON users USING GIN (first_name gin_trgm_ops, second_name gin_trgm_ops);

COMMIT;