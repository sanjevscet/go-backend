ALTER TABLE "users" ALTER COLUMN password TYPE bytea USING password::bytea;

