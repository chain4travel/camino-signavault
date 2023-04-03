CREATE TABLE multisig_tx
(
    id              CHAR(64)        NOT NULL,
    unsigned_tx     VARCHAR(2048)   NOT NULL,
    alias           VARCHAR(255)    NOT NULL,
    threshold       INT             NOT NULL,
    transaction_id  VARCHAR(56)     NULL,
    output_owners   VARCHAR(255)    NOT NULL,
    metadata        VARCHAR(255)    NOT NULL,
    created_at      DATETIME        NOT NULL,
    PRIMARY KEY (id)
);

CREATE INDEX idx_multisig_tx_alias ON multisig_tx (alias);
CREATE UNIQUE INDEX idx_multisig_tx_id ON multisig_tx (id);

CREATE TABLE multisig_tx_owners
(
    multisig_tx_id CHAR(64)         NOT NULL,
    address        CHAR(51)         NOT NULL,
    signature      VARCHAR(255)     NULL,
    is_signer      BOOLEAN          NOT NULL DEFAULT FALSE,
    created_at     DATETIME         NOT NULL,
    FOREIGN KEY (multisig_tx_id) REFERENCES multisig_tx (id),
    PRIMARY KEY (multisig_tx_id, address)
);