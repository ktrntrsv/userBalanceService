CREATE TABLE account
(
    id      uuid PRIMARY KEY,
    balance float check ( balance >= 0 ) NOT NULL
);

CREATE TABLE transaction
(
    id              uuid PRIMARY KEY,
    amount          float check ( amount >= 0 ) NOT NULL,
    account_to_id   uuid                        NOT NULL REFERENCES account (id),
    account_from_id uuid                        NOT NULL REFERENCES account (id),
    status          int4                        NOT NULL,
    service_id      varchar(50)
);

INSERT INTO account(id, balance)
VALUES ('123e4567-e89b-12d3-a456-426614174000', 0);

INSERT INTO account(id, balance)
VALUES ('123e4567-e89b-12d3-a456-426614174001', 0);

INSERT INTO account(id, balance)
VALUES ('123e4567-e89b-12d3-a456-426614174002', 0)