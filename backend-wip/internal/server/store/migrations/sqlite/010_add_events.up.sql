-- Add category and value columns to events table
-- (events table already exists from 001_initial)
ALTER TABLE events ADD COLUMN category TEXT;
ALTER TABLE events ADD COLUMN value TEXT;
