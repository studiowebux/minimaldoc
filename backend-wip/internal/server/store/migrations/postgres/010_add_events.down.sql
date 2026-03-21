-- Remove category and value columns from events table
ALTER TABLE events DROP COLUMN IF EXISTS category;
ALTER TABLE events DROP COLUMN IF EXISTS value;
