SET foreign_key_checks = 0;
ALTER TABLE multisig_tx MODIFY id char(64);
ALTER TABLE multisig_tx_owners MODIFY multisig_tx_id char(64);
SET foreign_key_checks = 1;