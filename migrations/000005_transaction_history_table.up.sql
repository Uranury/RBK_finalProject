CREATE TABLE IF NOT EXISTS transaction_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    type VARCHAR(20) NOT NULL CHECK (type IN ('withdraw', 'deposit')),
    balance_before DECIMAL(12,2) NOT NULL,
    balance_after DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transaction_history_user_id ON transaction_history(user_id);
CREATE INDEX idx_transaction_history_type ON transaction_history(type);
CREATE INDEX idx_transaction_history_created_at ON transaction_history(created_at DESC);
