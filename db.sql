CREATE SCHEMA IF NOT EXISTS rinha;

CREATE TABLE IF NOT EXISTS rinha.party (
    party_id BIGSERIAL PRIMARY KEY,
    "limit" BIGINT,
    balance BIGINT
    CHECK (balance > -"limit")
);

CREATE TABLE IF NOT EXISTS rinha.transactions (
    transaction_id BIGSERIAL PRIMARY KEY,
    value BIGINT,
    type VARCHAR(1),
    description VARCHAR(10),
    party_id BIGINT,
    created_at timestamp,
    CONSTRAINT party_fk
        FOREIGN KEY(party_id)
            REFERENCES rinha.party(party_id)
);

INSERT INTO rinha.party ("limit", balance) VALUES
(100000, 0),
(80000, 0),
(1000000, 0),
(10000000, 0),
(500000, 0);
