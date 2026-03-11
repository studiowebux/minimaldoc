-- SQLite: Author role validation handled at application level
-- No schema change needed for SQLite as CHECK constraint modification
-- requires table recreation. Application validates role on insert/update.

-- This migration exists for version tracking consistency with PostgreSQL.
SELECT 1;
