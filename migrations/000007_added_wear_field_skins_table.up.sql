-- Add wear field to skins table
ALTER TABLE skins ADD COLUMN wear VARCHAR(20) NOT NULL DEFAULT 'Field-Tested';

-- Create enum constraint for wear field based on CS:GO wear levels
ALTER TABLE skins ADD CONSTRAINT chk_wear_type CHECK (
    wear IN (
        'Factory New',
        'Minimal Wear', 
        'Field-Tested',
        'Well-Worn',
        'Battle-Scarred'
    )
);

-- Create index on wear field for better query performance
CREATE INDEX idx_skins_wear ON skins(wear);

-- Create a function to automatically calculate wear based on condition
CREATE OR REPLACE FUNCTION calculate_wear_from_condition(condition_value DECIMAL(10,8))
RETURNS VARCHAR(20) AS $$
BEGIN
    CASE 
        WHEN condition_value >= 0.00 AND condition_value <= 0.07 THEN
            RETURN 'Factory New';
        WHEN condition_value > 0.07 AND condition_value <= 0.15 THEN
            RETURN 'Minimal Wear';
        WHEN condition_value > 0.15 AND condition_value <= 0.38 THEN
            RETURN 'Field-Tested';
        WHEN condition_value > 0.38 AND condition_value <= 0.45 THEN
            RETURN 'Well-Worn';
        WHEN condition_value > 0.45 AND condition_value <= 1.00 THEN
            RETURN 'Battle-Scarred';
        ELSE
            RETURN 'Field-Tested'; -- Default fallback
    END CASE;
END;
$$ LANGUAGE plpgsql;

-- Update existing skins to have the correct wear based on their condition
UPDATE skins SET wear = calculate_wear_from_condition(condition);

-- Create a trigger to automatically update wear when condition changes
CREATE OR REPLACE FUNCTION update_wear_on_condition_change()
RETURNS TRIGGER AS $$
BEGIN
    NEW.wear = calculate_wear_from_condition(NEW.condition);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_wear
    BEFORE INSERT OR UPDATE ON skins
    FOR EACH ROW
    EXECUTE FUNCTION update_wear_on_condition_change();
