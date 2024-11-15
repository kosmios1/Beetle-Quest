CREATE TABLE transactions (
    transaction_id UUID PRIMARY KEY,
    transaction_type INTEGER NOT NULL,

    user_id UUID NOT NULL,

    amount BIGINT NOT NULL,
    date_time TIMESTAMP NOT NULL,

    event_type INTEGER NOT NULL,
    event_id UUID NOT NULL
);
