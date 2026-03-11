-- Add missing columns to events table for custom event tracking
-- The events table already exists from 001_initial.up.sql with: id, site_id, name, path, properties, session_hash, created_at
-- We add category and value columns for simpler event data

ALTER TABLE events ADD COLUMN IF NOT EXISTS category TEXT;
ALTER TABLE events ADD COLUMN IF NOT EXISTS value TEXT;
