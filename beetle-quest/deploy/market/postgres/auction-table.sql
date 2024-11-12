CREATE TABLE bids (
    auction_id UUID,
    user_id UUID,
    amount_spend BIGINT,
    timestamp TIMESTAMP
);

CREATE TABLE auctions (
    auction_id UUID PRIMARY KEY,
    owner_id UUID NOT NULL,
    gacha_id UUID NOT NULL,

    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,

    winner_id UUID
);
