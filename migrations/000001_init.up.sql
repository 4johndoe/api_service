CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE currency
(
    id     SERIAL PRIMARY KEY,
    name   TEXT,
    symbol TEXT
);

CREATE TABLE category
(
    id   SERIAL PRIMARY KEY,
    name TEXT
);

CREATE TABLE product
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          TEXT NOT NULL,
    description   TEXT NOT NULL,
    image_id      UUID,
    price         BIGINT,
    currency_id   INT REFERENCES currency (id),
    rating        INT,
    category_id   INT REFERENCES category (id),
    specification JSONB,
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ
    CONSTRAINT positive_price CHECK (price > 0)
    CONSTRAINT valid_rating CHECK (rating <= 5)
);

-- DATA --

INSERT INTO currency (name, symbol)
VALUES ('USD', '$');
INSERT INTO currency (name, symbol)
VALUES ('EUR', '€');

INSERT INTO category (name)
VALUES ('купоны');
INSERT INTO category (name)
VALUES ('цифровые билеты');

INSERT INTO product (name, description, image_id, price, currency_id, rating, category_id, specification, created_at, updated_at)
VALUES ('Cellphone', 'Mobile phone some text', null, 100, 1, 1, 1, null, now(), now())