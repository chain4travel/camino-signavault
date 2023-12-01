SET foreign_key_checks = 0;
ALTER TABLE multisig_tx_owners MODIFY multisig_tx_id char(84);
ALTER TABLE multisig_tx MODIFY id char(84);
SET foreign_key_checks = 1;