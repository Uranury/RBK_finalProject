ALTER TABLE transaction_history
ADD CONSTRAINT transaction_history_amount_check CHECK (amount > 0);