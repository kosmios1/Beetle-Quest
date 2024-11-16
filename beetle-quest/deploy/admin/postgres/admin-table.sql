CREATE TABLE admins (
    admin_id UUID PRIMARY KEY,
    password_hash BYTEA NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    opt_secret BYTEA NOT NULL
);


INSERT INTO admins (admin_id, password_hash, email, opt_secret) VALUES
-- password = admin
('09087f45-5209-4efa-85bd-761562a6df53', decode('243261243130247370373732344b6d544a302e4f347862557176514d754c5330464a79684e4355736c6e59787757685a6668386a7739704430644457', 'hex'), 'admin@admin.com', decode('393dafd43eb2402579de62775c138347474c10f6', 'hex'));
