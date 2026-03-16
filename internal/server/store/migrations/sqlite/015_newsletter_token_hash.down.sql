-- SQLite cannot DROP COLUMN before 3.35; recreate table if needed.
-- For rollback, just clear the new column's data.
UPDATE subscribers SET verify_token_expires_at = NULL;
