CREATE TABLE multisig_tx
(
    id             INT          NOT NULL AUTO_INCREMENT,
    alias          CHAR(33)     NOT NULL,
    threshold      INT          NOT NULL,
    unsignedTx     VARCHAR(255) NOT NULL,
    transaction_id CHAR(49)     NULL,
    created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE INDEX idx_multisig_tx_alias ON multisig_tx (alias);

CREATE TABLE multisig_tx_signers
(
    id             INT          NOT NULL AUTO_INCREMENT,
    multisig_tx_id INT          NOT NULL,
    address        CHAR(33)     NOT NULL,
    signature      VARCHAR(255) NOT NULL,
    created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (multisig_tx_id) REFERENCES multisig_tx (id)
);