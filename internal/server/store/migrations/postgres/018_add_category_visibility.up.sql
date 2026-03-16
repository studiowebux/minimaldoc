-- Add visibility column to forum_categories for access control.
ALTER TABLE forum_categories ADD COLUMN visibility TEXT DEFAULT 'public' CHECK (visibility IN ('public', 'members_only', 'staff_only'));
