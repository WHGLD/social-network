-- +migrate Down

BEGIN;

DROP TYPE gender

DROP TABLE IF EXISTS users CASCADE;

COMMIT;