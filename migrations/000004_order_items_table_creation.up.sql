CREATE TABLE IF NOT EXISTS order_items (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
       skin_id UUID NOT NULL REFERENCES skins(id) ON DELETE RESTRICT,
       price DECIMAL(12,2) NOT NULL CHECK (price > 0),
       created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_skin_id ON order_items(skin_id);