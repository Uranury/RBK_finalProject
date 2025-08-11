-- Add gun field to skins table
ALTER TABLE skins ADD COLUMN gun VARCHAR(50) NOT NULL DEFAULT 'AK-47';

-- Create enum constraint for gun field
ALTER TABLE skins ADD CONSTRAINT chk_gun_type CHECK (
    gun IN (
        -- Pistols
        'AK-47', 'M4A4', 'M4A1-S', 'Desert Eagle', 'USP-S', 'Glock-18', 'P250', 'Tec-9', 'CZ75-Auto',
        -- Rifles
        'AWP', 'SSG 08', 'SCAR-20', 'G3SG1',
        -- SMGs
        'MP9', 'MAC-10', 'MP7', 'P90', 'UMP-45', 'PP-Bizon',
        -- Shotguns
        'Nova', 'XM1014', 'MAG-7', 'Sawed-Off',
        -- Machine Guns
        'M249', 'Negev',
        -- Knives
        'Karambit', 'Butterfly Knife', 'M9 Bayonet', 'Bayonet', 'Flip Knife', 'Gut Knife', 
        'Huntsman Knife', 'Shadow Daggers',
        -- Other
        'Falchion Knife', 'Bowie Knife', 'Navaja Knife', 'Stiletto Knife', 'Ursus Knife', 
        'Nomad Knife', 'Paracord Knife', 'Survival Knife', 'Classic Knife'
    )
);

-- Create index on gun field for better query performance
CREATE INDEX idx_skins_gun ON skins(gun);
