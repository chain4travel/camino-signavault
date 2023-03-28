/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package dao

import (
	"database/sql"
	"log"

	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/model"
)

type MultisigTxDao interface {
	CreateMultisigTx(id string, alias string, threshold int, unsignedTx string, creator string, signature string, outputOwners string, metadata string, owners []string) (string, error)
	GetMultisigTx(id string, alias string, owner string) (*[]model.MultisigTx, error)
	UpdateTransactionId(id string, transactionId string) (bool, error)
	AddSigner(id string, signature string, signerAddress string) (bool, error)
	PendingAliasExists(alias string) (bool, error)
}
type multisigTxDao struct {
	db *db.Db
}

func NewMultisigTxDao(db *db.Db) MultisigTxDao {
	return &multisigTxDao{
		db: db,
	}
}

func (d *multisigTxDao) PendingAliasExists(alias string) (bool, error) {
	query := "SELECT count(id) " +
		"FROM multisig_tx " +
		"WHERE alias = ? AND transaction_id IS NULL"
	rows, err := d.db.Query(query, alias)
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

func (d *multisigTxDao) CreateMultisigTx(id string, alias string, threshold int, unsignedTx string, creator string, signature string, outputOwners string, metadata string, owners []string) (string, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return "", err
	}

	stmt, err := tx.Prepare("INSERT INTO multisig_tx (id, alias, threshold, unsigned_tx, output_owners, metadata) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return "", err
	}
	_, err = stmt.Exec(id, alias, threshold, unsignedTx, outputOwners, metadata)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Execute statement failed: %v, unable to rollback: %v", err, rollbackErr)
		}
		log.Print(err)
		return "", err
	}

	for _, owner := range owners {
		isSigner := false
		ownerSignature := ""
		if owner == creator {
			isSigner = true
			ownerSignature = signature
		}

		stmt, err := tx.Prepare("INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer) VALUES (?, ?, ?, ?)")
		if err != nil {
			return "", err
		}
		_, err = stmt.Exec(id, owner, ownerSignature, isSigner)
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

	return id, nil
}

func (d *multisigTxDao) GetMultisigTx(id string, alias string, owner string) (*[]model.MultisigTx, error) {
	var err error

	var query string
	var rows *sql.Rows
	if owner == "" {
		query = "SELECT tx.id, " +
			"tx.alias, " +
			"tx.threshold, " +
			"tx.transaction_id, " +
			"tx.unsigned_tx, " +
			"tx.output_owners," +
			"tx.metadata," +
			"owners.multisig_tx_id, " +
			"owners.address, " +
			"owners.signature, " +
			"owners.is_signer " +
			"FROM multisig_tx AS tx " +
			"LEFT JOIN multisig_tx_owners AS owners ON tx.id = owners.multisig_tx_id " +
			"WHERE (tx.alias=? OR ?='') AND (tx.id=? OR ?='') AND tx.transaction_id IS NULL " +
			"ORDER BY tx.created_at ASC"
		rows, err = d.db.Query(query, alias, alias, id, id)
	} else {
		query = "SELECT tx.id, " +
			"tx.alias, " +
			"tx.threshold, " +
			"tx.transaction_id, " +
			"tx.unsigned_tx, " +
			"tx.output_owners," +
			"tx.metadata," +
			"owners.multisig_tx_id, " +
			"owners.address, " +
			"owners.signature, " +
			"owners.is_signer, " +
			"owners2.address " +
			"FROM multisig_tx AS tx " +
			"LEFT JOIN multisig_tx_owners AS owners ON tx.id = owners.multisig_tx_id " +
			"JOIN multisig_tx_owners AS owners2 ON tx.id = owners2.multisig_tx_id " +
			"WHERE (tx.alias=? OR ?='') AND (tx.id=? OR ?='') AND (owners2.address = ? OR ?='') AND tx.transaction_id IS NULL " +
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
			txOutputOwners    string
			txMetadata        string
			ownerMultisigTxId string
			ownerAddress      sql.NullString
			ownerSignature    sql.NullString
			ownerIsSigner     sql.NullBool
			ownerAddress2     sql.NullString
		)

		var err error
		if owner == "" {
			err = rows.Scan(&txId, &txAlias, &txThreshold, &txTransactionId, &txUnsignedTx, &txOutputOwners, &txMetadata, &ownerMultisigTxId, &ownerAddress, &ownerSignature, &ownerIsSigner)
		} else {
			err = rows.Scan(&txId, &txAlias, &txThreshold, &txTransactionId, &txUnsignedTx, &txOutputOwners, &txMetadata, &ownerMultisigTxId, &ownerAddress, &ownerSignature, &ownerIsSigner, &ownerAddress2)
		}
		if err != nil {
			log.Fatal(err)
		}

		var tx model.MultisigTx
		if _, ok := multiSigTx[txId]; ok {
			tx = multiSigTx[txId]
		} else {
			tx = model.MultisigTx{
				Id:            txId,
				Alias:         txAlias,
				Threshold:     txThreshold,
				OutputOwners:  txOutputOwners,
				Metadata:      txMetadata,
				TransactionId: txTransactionId.String,
				UnsignedTx:    txUnsignedTx,
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
