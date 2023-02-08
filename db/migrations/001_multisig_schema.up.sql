CREATE TABLE multisig_tx
(
    id             INT          NOT NULL AUTO_INCREMENT,
    alias          CHAR(51)     NOT NULL,
    threshold      INT          NOT NULL,
    unsigned_tx    VARCHAR(255) NOT NULL,
    transaction_id CHAR(51)     NULL,
    created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT uc_multisig_tx_alias UNIQUE (alias)
);

CREATE INDEX idx_multisig_tx_alias ON multisig_tx (alias);

CREATE TABLE multisig_tx_signers
(
    id             INT          NOT NULL AUTO_INCREMENT,
    multisig_tx_id INT          NOT NULL,
    address        CHAR(51)     NOT NULL,
    signature      VARCHAR(255) NOT NULL,
    created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (multisig_tx_id) REFERENCES multisig_tx (id),
    CONSTRAINT uc_multisig_tx_id_address UNIQUE (multisig_tx_id, address)
);