ALTER TABLE multisig_tx MODIFY unsigned_tx VARCHAR(4096) CHARACTER SET utf8mb4;
ALTER TABLE multisig_tx MODIFY output_owners VARCHAR(512) CHARACTER SET utf8mb4;