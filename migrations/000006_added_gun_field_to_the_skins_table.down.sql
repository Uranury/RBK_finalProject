-- Remove index on gun field
DROP INDEX IF EXISTS idx_skins_gun;

-- Remove gun field from skins table
ALTER TABLE skins DROP COLUMN IF EXISTS gun;
