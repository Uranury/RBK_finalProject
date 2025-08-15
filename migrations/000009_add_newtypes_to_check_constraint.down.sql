ALTER TABLE transaction_history
    DROP CONSTRAINT transaction_history_type_check,
    ADD CONSTRAINT transaction_history_type_check
        CHECK (type IN ('withdraw', 'deposit'));
