BEGIN;

ALTER TABLE discount_histories ALTER COLUMN started_at SET NOT NULL;
ALTER TABLE discount_histories ALTER COLUMN ended_at SET NOT NULL;

COMMIT;