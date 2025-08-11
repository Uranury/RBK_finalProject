-- Drop the trigger first
DROP TRIGGER IF EXISTS trigger_update_wear ON skins;

-- Drop the functions
DROP FUNCTION IF EXISTS update_wear_on_condition_change();
DROP FUNCTION IF EXISTS calculate_wear_from_condition(DECIMAL);

-- Remove index on wear field
DROP INDEX IF EXISTS idx_skins_wear;

-- Remove wear field from skins table
ALTER TABLE skins DROP COLUMN IF EXISTS wear;
