CREATE TABLE deposit_offer_sigs
(
    deposit_offer_id CHAR(64) NOT NULL,
    address          CHAR(51) NOT NULL,
    signature        VARCHAR(255) NULL,
    PRIMARY KEY (deposit_offer_id, address)
);

CREATE INDEX whitelisted_addr ON deposit_offer_sigs (address);