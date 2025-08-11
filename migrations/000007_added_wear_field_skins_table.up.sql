ALTER TABLE skins
ADD COLUMN wear TEXT NOT NULL, 
ADD CONSTRAINT chk_wear_type CHECK (wear IN ('Factory New', 'Minimal Wear', 'Field-Tested', 'Well-Worn', 'Battle-Scarred'));
