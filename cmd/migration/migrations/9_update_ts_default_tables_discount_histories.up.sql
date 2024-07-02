BEGIN;

ALTER TABLE discount_histories ALTER COLUMN started_at DROP DEFAULT;
ALTER TABLE discount_histories ALTER COLUMN started_at DROP NOT NULL;
ALTER TABLE discount_histories ALTER COLUMN ended_at DROP DEFAULT;
ALTER TABLE discount_histories ALTER COLUMN ended_at DROP NOT NULL;

COMMIT;