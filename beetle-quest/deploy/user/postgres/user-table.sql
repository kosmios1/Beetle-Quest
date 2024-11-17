CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    currency BIGINT NOT NULL,
    password_hash BYTEA NOT NULL
);

INSERT INTO users (user_id, username, email, currency, password_hash) VALUES
-- password = admin
('09087f45-5209-4efa-85bd-761562a6df53', 'admin', 'admin@admin.com', 10000, decode('243261243130247370373732344b6d544a302e4f347862557176514d754c5330464a79684e4355736c6e59787757685a6668386a7739704430644457', 'hex'));
