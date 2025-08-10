CREATE TABLE IF NOT EXISTS skins (
     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     owner_id UUID REFERENCES users(id) ON DELETE SET NULL,
     name VARCHAR(255) NOT NULL,
     rarity VARCHAR(100) NOT NULL,
     condition DECIMAL(10,8) NOT NULL CHECK (condition >= 0.0 AND condition <= 1.0),
     price DECIMAL(12,2) NOT NULL CHECK (price >= 0),
     image VARCHAR(500),
     available BOOLEAN DEFAULT true,
     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
     updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_skins_available ON skins(available) WHERE available = true;
CREATE INDEX idx_skins_owner ON skins(owner_id);
CREATE INDEX idx_skins_rarity ON skins(rarity);
CREATE INDEX idx_skins_price ON skins(price);