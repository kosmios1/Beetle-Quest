CREATE TABLE user_gacha (
    user_id UUID NOT NULL,
    gacha_id UUID NOT NULL,
    PRIMARY KEY (user_id, gacha_id) UNIQUE,
);
