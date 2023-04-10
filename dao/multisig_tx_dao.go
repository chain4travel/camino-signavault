/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package dao

import (
	"database/sql"
	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/model"
	"log"
	"time"
)

type MultisigTxDao interface {
	CreateMultisigTx(multisig *model.MultisigTx) (string, error)
	GetMultisigTx(id string, alias string, owner string) (*[]model.MultisigTx, error)
	UpdateTransactionId(id string, transactionId string) (bool, error)
	AddSigner(id string, signature string, signerAddress string) (bool, error)
	PendingAliasExists(alias string, chainId string) (bool, error)
}
type multisigTxDao struct {
	db *db.Db
}

func NewMultisigTxDao(db *db.Db) MultisigTxDao {
	return &multisigTxDao{
		db: db,
	}
}

func (d *multisigTxDao) PendingAliasExists(alias string, chainId string) (bool, error) {
	query := "SELECT count(id) " +
		"FROM multisig_tx " +
		"WHERE alias = ? AND chain_id = ? AND transaction_id IS NULL AND (expires_at > UTC_TIMESTAMP() OR expires_at IS NULL)"
	rows, err := d.db.Query(query, alias, chainId)
	if err != nil {
		return false, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Print(err)
		}
	}(rows)

	rows.Next()
	var count int
	err = rows.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *multisigTxDao) CreateMultisigTx(multisig *model.MultisigTx) (string, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return "", err
	}

	stmt, err := tx.Prepare("INSERT INTO multisig_tx (id, alias, threshold, chain_id, unsigned_tx, output_owners, metadata, parent_transaction, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return "", err
	}
	now := time.Now().UTC()
	_, err = stmt.Exec(multisig.Id, multisig.Alias, multisig.Threshold, multisig.ChainId, multisig.UnsignedTx, multisig.OutputOwners, multisig.Metadata, multisig.ParentTransaction, multisig.Expiration, now)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Execute statement failed: %v, unable to rollback: %v", err, rollbackErr)
		}
		log.Print(err)
		return "", err
	}

	owners := multisig.Owners
	for _, owner := range owners {
		stmt, err := tx.Prepare("INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer, created_at) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			return "", err
		}
		_, err = stmt.Exec(multisig.Id, owner.Address, owner.Signature, owner.Signature != "", now)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Execute statement failed: %v, unable to rollback: %v", err, rollbackErr)
			}
			log.Print(err)
			return "", err
		}

	}
	err = tx.Commit()
	if err != nil {
		log.Printf("Commit failed: %v", err)
		return "", err
	}

	return multisig.Id, nil
}

func (d *multisigTxDao) GetMultisigTx(id string, alias string, owner string) (*[]model.MultisigTx, error) {
	var err error

	var query string
	var rows *sql.Rows
	if owner == "" {
		query = "SELECT tx.id, " +
			"tx.alias, " +
			"tx.threshold, " +
			"tx.chain_id, " +
			"tx.transaction_id, " +
			"tx.unsigned_tx, " +
			"tx.output_owners," +
			"tx.metadata," +
			"tx.parent_transaction," +
			"tx.expires_at," +
			"tx.created_at," +
			"owners.multisig_tx_id, " +
			"owners.address, " +
			"owners.signature, " +
			"owners.is_signer " +
			"FROM multisig_tx AS tx " +
			"LEFT JOIN multisig_tx_owners AS owners ON tx.id = owners.multisig_tx_id " +
			"WHERE (tx.alias=? OR ?='') AND (tx.id=? OR ?='') AND tx.transaction_id IS NULL AND (tx.expires_at > UTC_TIMESTAMP() OR tx.expires_at IS NULL)" +
			"ORDER BY tx.created_at ASC"
		rows, err = d.db.Query(query, alias, alias, id, id)
	} else {
		query = "SELECT tx.id, " +
			"tx.alias, " +
			"tx.threshold, " +
			"tx.chain_id, " +
			"tx.transaction_id, " +
			"tx.unsigned_tx, " +
			"tx.output_owners," +
			"tx.metadata," +
			"tx.parent_transaction," +
			"tx.expires_at," +
			"tx.created_at," +
			"owners.multisig_tx_id, " +
			"owners.address, " +
			"owners.signature, " +
			"owners.is_signer, " +
			"owners2.address " +
			"FROM multisig_tx AS tx " +
			"LEFT JOIN multisig_tx_owners AS owners ON tx.id = owners.multisig_tx_id " +
			"JOIN multisig_tx_owners AS owners2 ON tx.id = owners2.multisig_tx_id " +
			"WHERE (tx.alias=? OR ?='') AND (tx.id=? OR ?='') AND (owners2.address = ? OR ?='') AND tx.transaction_id IS NULL AND (tx.expires_at > UTC_TIMESTAMP() OR tx.expires_at IS NULL)" +
			"ORDER BY tx.created_at ASC"
		rows, err = d.db.Query(query, alias, alias, id, id, owner, owner)
	}

	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Print(err)
		}
	}(rows)

	var result []model.MultisigTx
	multiSigTx := make(map[string]model.MultisigTx)

	for rows.Next() {
		var (
			txId              string
			txAlias           string
			txThreshold       int8
			txTransactionId   sql.NullString
			txUnsignedTx      string
			txChainId         string
			txOutputOwners    string
			txMetadata        string
			txParentTx        sql.NullString
			txExpiresAt       sql.NullTime
			txCreatedAt       time.Time
			ownerMultisigTxId string
			ownerAddress      sql.NullString
			ownerSignature    sql.NullString
			ownerIsSigner     sql.NullBool
			ownerAddress2     sql.NullString
		)

		var err error
		if owner == "" {
			err = rows.Scan(&txId, &txAlias, &txThreshold, &txChainId, &txTransactionId, &txUnsignedTx, &txOutputOwners,
				&txMetadata, &txParentTx, &txExpiresAt, &txCreatedAt, &ownerMultisigTxId, &ownerAddress, &ownerSignature, &ownerIsSigner)
		} else {
			err = rows.Scan(&txId, &txAlias, &txThreshold, &txChainId, &txTransactionId, &txUnsignedTx, &txOutputOwners,
				&txMetadata, &txParentTx, &txExpiresAt, &txCreatedAt, &ownerMultisigTxId, &ownerAddress, &ownerSignature, &ownerIsSigner, &ownerAddress2)
		}
		if err != nil {
			log.Fatal(err)
		}

		var tx model.MultisigTx
		if _, ok := multiSigTx[txId]; ok {
			tx = multiSigTx[txId]
		} else {

			var expiration *time.Time
			if txExpiresAt.Valid {
				t := txExpiresAt.Time.UTC()
				expiration = &t
			}
			t := txCreatedAt.UTC()
			created := &t

			tx = model.MultisigTx{
				Id:                txId,
				UnsignedTx:        txUnsignedTx,
				Alias:             txAlias,
				Threshold:         txThreshold,
				ChainId:           txChainId,
				TransactionId:     txTransactionId.String,
				OutputOwners:      txOutputOwners,
				Metadata:          txMetadata,
				ParentTransaction: txParentTx.String,
				Expiration:        expiration,
				Timestamp:         created,
			}
		}

		owners := tx.Owners
		if owners == nil {
			owners = []model.MultisigTxOwner{}
		}
		// add owner
		owner := model.MultisigTxOwner{
			MultisigTxId: ownerMultisigTxId,
			Address:      ownerAddress.String,
			Signature:    ownerSignature.String,
		}
		owners = append(owners, owner)
		tx.Owners = owners

		multiSigTx[txId] = tx

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	// convert map to slice
	for _, tx := range multiSigTx {
		result = append(result, tx)
	}
	if result == nil {
		return nil, nil
	}
	return &result, nil
}

func (d *multisigTxDao) UpdateTransactionId(id string, transactionId string) (bool, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE multisig_tx SET transaction_id = ? WHERE id = ? AND transaction_id IS NULL")
	if err != nil {
		return false, err
	}
	_, err = stmt.Exec(transactionId, id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Execute statement failed: %v, unable to rollback: %v", err, rollbackErr)
		}
		log.Print(err)
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Commit failed: %v", err)
		return false, err
	}

	return true, nil
}

func (d *multisigTxDao) AddSigner(id string, signature string, signerAddress string) (bool, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE multisig_tx_owners SET signature = ?, is_signer = ? WHERE multisig_tx_id = ? AND address = ?")
	if err != nil {
		return false, err
	}
	_, err = stmt.Exec(signature, true, id, signerAddress)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Execute statement failed: %v, unable to rollback: %v", err, rollbackErr)
		}
		log.Print(err)
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Commit failed: %v", err)
		return false, err
	}

	return true, nil
}
