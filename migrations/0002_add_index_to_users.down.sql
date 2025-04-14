-- +migrate Down

BEGIN;

DROP INDEX IF EXISTS users_names_ind;

COMMIT;