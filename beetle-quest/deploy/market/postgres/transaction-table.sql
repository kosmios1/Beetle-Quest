CREATE TABLE transactions (
    transaction_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,

    transaction_type INTEGER NOT NULL,

    amount BIGINT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    event_type INTEGER NOT NULL,
    event_id UUID NOT NULL
);
