INSERT INTO multisig_tx (id, unsigned_tx, alias, threshold, metadata, output_owners, created_at) VALUES ('1', 'unsigned_tx', 'alias', 2, 'metadata', 'output_owners', NOW());
INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer, created_at) VALUES ('1', 'address', 'signature', true, NOW());

INSERT INTO multisig_tx (id, unsigned_tx, alias, threshold, metadata, output_owners, transaction_id, created_at) VALUES ('2', 'unsigned_tx_2', 'alias_2', 3, 'metadata_2', 'output_owners_2', 'transaction_id_2', NOW());
INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer, created_at) VALUES ('2', 'address1', 'signature1', true, NOW());
INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer, created_at) VALUES ('2', 'address2', 'signature2', true, NOW());
INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer, created_at) VALUES ('2', 'address3', 'signature3', true, NOW());

INSERT INTO multisig_tx (id, unsigned_tx, alias, threshold, metadata, output_owners, created_at) VALUES ('3', 'unsigned_tx_3', 'alias_3', 2, 'metadata_3','output_owners_3', NOW());
INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer, created_at) VALUES ('3', 'address1', 'signature1', true, NOW());
INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer, created_at) VALUES ('3', 'address2', 'signature2', true, NOW());