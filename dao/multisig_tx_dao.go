package dao

import (
	"database/sql"
	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/model"
	"log"
)

type MultisigTxDaoInterface interface {
	CreateMultisigTx(alias string, threshold int, unsignedTx string, creator string, signature string, owners []string) (int64, error)
	GetMultisigTx(id int64, alias string, owner string) (*[]model.MultisigTx, error)
	UpdateTransactionId(id int64, transactionId string) (bool, error)
	AddSigner(id int64, signature string, signerAddress string) (bool, error)
}
type MultisigTxDao struct {
	db *db.Db
}

func NewMultisigTxDao(db *db.Db) MultisigTxDaoInterface {
	return &MultisigTxDao{
		db: db,
	}
}

func (d *MultisigTxDao) CreateMultisigTx(alias string, threshold int, unsignedTx string, creator string, signature string, owners []string) (int64, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return -1, err
	}

	stmt, err := tx.Prepare("INSERT INTO multisig_tx (alias, threshold, unsigned_tx) VALUES (?, ?, ?)")
	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec(alias, threshold, unsignedTx)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Execute statement failed: %v, unable to rollback: %v", err, rollbackErr)
		}
		log.Print(err)
		return -1, err
	}
	txId, _ := res.LastInsertId()

	for _, owner := range owners {
		isSigner := false
		ownerSignature := ""
		if owner == creator {
			isSigner = true
			ownerSignature = signature
		}

		stmt, err := tx.Prepare("INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer) VALUES (?, ?, ?, ?)")
		if err != nil {
			return -1, err
		}
		_, err = stmt.Exec(txId, owner, ownerSignature, isSigner)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Execute statement failed: %v, unable to rollback: %v", err, rollbackErr)
			}
			log.Print(err)
			return -1, err
		}

	}
	err = tx.Commit()
	if err != nil {
		log.Printf("Commit failed: %v", err)
		return -1, err
	}

	return txId, nil
}

func (d *MultisigTxDao) GetMultisigTx(id int64, alias string, owner string) (*[]model.MultisigTx, error) {
	var err error

	var query string
	var rows *sql.Rows
	if owner == "" {
		query = "SELECT tx.id, " +
			"tx.alias, " +
			"tx.threshold, " +
			"tx.transaction_id, " +
			"tx.unsigned_tx, " +
			"owners.multisig_tx_id, " +
			"owners.id, " +
			"owners.address, " +
			"owners.signature, " +
			"owners.is_signer " +
			"FROM multisig_tx AS tx " +
			"LEFT JOIN multisig_tx_owners AS owners ON tx.id = owners.multisig_tx_id " +
			"WHERE (tx.alias=? OR ?='') AND (tx.id=? OR ?=-1) AND tx.transaction_id IS NULL " +
			"ORDER BY tx.created_at ASC"
		rows, err = d.db.Query(query, alias, alias, id, id)
	} else {
		query = "SELECT tx.id, " +
			"tx.alias, " +
			"tx.threshold, " +
			"tx.transaction_id, " +
			"tx.unsigned_tx, " +
			"owners.multisig_tx_id, " +
			"owners.id, " +
			"owners.address, " +
			"owners.signature, " +
			"owners.is_signer, " +
			"owners2.address " +
			"FROM multisig_tx AS tx " +
			"LEFT JOIN multisig_tx_owners AS owners ON tx.id = owners.multisig_tx_id " +
			"JOIN multisig_tx_owners AS owners2 ON tx.id = owners2.multisig_tx_id " +
			"WHERE (tx.alias=? OR ?='') AND (tx.id=? OR ?=-1) AND (owners2.address = ? OR ?='') AND tx.transaction_id IS NULL " +
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
	multiSigTx := make(map[int64]model.MultisigTx)

	for rows.Next() {
		var (
			txId              int64
			txAlias           string
			txThreshold       int8
			txTransactionId   sql.NullString
			txUnsignedTx      string
			ownerMultisigTxId sql.NullInt64
			ownerId           sql.NullInt64
			ownerAddress      sql.NullString
			ownerSignature    sql.NullString
			ownerIsSigner     sql.NullBool
			ownerAddress2     sql.NullString
		)

		var err error
		if owner == "" {
			err = rows.Scan(&txId, &txAlias, &txThreshold, &txTransactionId, &txUnsignedTx, &ownerMultisigTxId, &ownerId, &ownerAddress, &ownerSignature, &ownerIsSigner)
		} else {
			err = rows.Scan(&txId, &txAlias, &txThreshold, &txTransactionId, &txUnsignedTx, &ownerMultisigTxId, &ownerId, &ownerAddress, &ownerSignature, &ownerIsSigner, &ownerAddress2)
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
			Id:           ownerId.Int64,
			MultisigTxId: ownerMultisigTxId.Int64,
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

func (d *MultisigTxDao) UpdateTransactionId(id int64, transactionId string) (bool, error) {
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

func (d *MultisigTxDao) AddSigner(id int64, signature string, signerAddress string) (bool, error) {
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
