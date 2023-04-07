INSERT INTO multisig_tx (id, unsigned_tx, alias, threshold, chain_id, metadata, output_owners, expires_at, created_at)
VALUES ('1', 'unsigned_tx', 'alias', 2, '11111111111111111111111111111111LpoYY', 'metadata', 'output_owners', NOW() + INTERVAL 1 YEAR, NOW());
INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer, created_at)
VALUES ('1', 'address', 'signature', true, NOW());

INSERT INTO multisig_tx (id, unsigned_tx, alias, threshold, chain_id, metadata, output_owners, transaction_id, expires_at, created_at)
VALUES ('2', 'unsigned_tx_2', 'alias_2', 3, '11111111111111111111111111111111LpoYY', 'metadata_2', 'output_owners_2', 'transaction_id_2', NOW() + INTERVAL 1 YEAR, NOW());
INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer, created_at)
VALUES ('2', 'address1', 'signature1', true, NOW());
INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer, created_at)
VALUES ('2', 'address2', 'signature2', true, NOW());
INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer, created_at)
VALUES ('2', 'address3', 'signature3', true, NOW());

INSERT INTO multisig_tx (id, unsigned_tx, alias, threshold, chain_id, metadata, output_owners, expires_at, created_at)
VALUES ('3', 'unsigned_tx_3', 'alias_3', 2, '11111111111111111111111111111111LpoYY', 'metadata_3','output_owners_3', NOW() + INTERVAL 1 YEAR, NOW());
INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer, created_at)
VALUES ('3', 'address1', 'signature1', true, NOW());
INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer, created_at)
VALUES ('3', 'address2', 'signature2', true, NOW());

INSERT INTO multisig_tx (id, unsigned_tx, alias, threshold, chain_id, metadata, output_owners, transaction_id, expires_at, created_at)
VALUES ('4', 'unsigned_tx_4', 'alias_3', 2, '11111111111111111111111111111111LpoYY', 'metadata_3','output_owners_3', 'transaction_id_3', NOW() + INTERVAL 1 YEAR, NOW());
INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer, created_at)
VALUES ('4', 'address1', 'signature1', true, NOW());

INSERT INTO multisig_tx (id, unsigned_tx, alias, threshold, chain_id, metadata, output_owners, expires_at, created_at)
VALUES ('5', 'unsigned_tx_5', 'alias_3', 2, 'jvYyfQTxGMJLuGWa55kdP2p2zSUYsQ5Raupu4TW34ZAUBAbtq', 'metadata_3','output_owners_3', NOW() + INTERVAL 1 YEAR, NOW());
INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer, created_at)
VALUES ('5', 'address1', 'signature1', true, NOW());

INSERT INTO multisig_tx (id, unsigned_tx, alias, threshold, chain_id, metadata, output_owners, transaction_id, expires_at, created_at)
VALUES ('6', 'unsigned_tx_6', 'alias_6', 2, 'jvYyfQTxGMJLuGWa55kdP2p2zSUYsQ5Raupu4TW34ZAUBAbtq', 'metadata_3','output_owners_3', 'transaction_id_6', NOW() + INTERVAL 1 YEAR, NOW());
INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer, created_at)
VALUES ('6', 'address1', 'signature1', true, NOW());