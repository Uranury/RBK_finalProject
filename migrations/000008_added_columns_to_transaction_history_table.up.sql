ALTER TABLE transaction_history
    ADD COLUMN skin_id UUID,
    ADD COLUMN order_id UUID,
    ADD COLUMN counterparty_id UUID,
    ADD COLUMN description TEXT;