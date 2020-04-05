CREATE TABLE IF NOT EXISTS account
(
    id       serial PRIMARY KEY,
    username VARCHAR(50) UNIQUE  NOT NULL,
    password VARCHAR(500)         NOT NULL,
    email    VARCHAR(355) UNIQUE NOT NULL
);