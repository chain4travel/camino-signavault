package service

import (
	"database/sql"
	"errors"
	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/model"
	"log"
)

type MultisigService struct {
	db db.Db
}

func NewMultisigService(db db.Db) *MultisigService {
	return &MultisigService{
		db: db,
	}
}

func (s *MultisigService) CreateMultisigTx(multisigTx *model.MultisigTx) (*model.MultisigTx, error) {
	var err error

	// check signers count is less than threshold
	signers := multisigTx.Signers
	if len(signers) >= multisigTx.Threshold {
		return nil, errors.New("signer count is more than threshold")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare("INSERT INTO multisig_tx (alias, threshold, unsigned_tx) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(multisigTx.Alias, multisigTx.Threshold, multisigTx.UnsignedTx)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Execute statement failed: %v, unable to rollback: %v", err, rollbackErr)
		}
		log.Print(err)
		return nil, err
	}
	txId, _ := res.LastInsertId()

	for _, owner := range multisigTx.Owners {
		isSigner := false
		signature := ""
		for _, signer := range multisigTx.Signers {
			if owner.Address == signer.Address {
				isSigner = true
				// check if signature is not empty
				if len(signer.Signature) == 0 {
					if rollbackErr := tx.Rollback(); rollbackErr != nil {
						log.Printf("Execute statement failed: %v, unable to rollback: %v", err, rollbackErr)
					}
					return nil, errors.New("signer signature is empty")
				}
				signature = signer.Signature
				break
			}
		}

		stmt, err := tx.Prepare("INSERT INTO multisig_tx_owners (multisig_tx_id, address, signature, is_signer) VALUES (?, ?, ?, ?)")
		if err != nil {
			return nil, err
		}
		_, err = stmt.Exec(txId, owner.Address, signature, isSigner)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Execute statement failed: %v, unable to rollback: %v", err, rollbackErr)
			}
			log.Print(err)
			return nil, err
		}

	}
	err = tx.Commit()
	if err != nil {
		log.Printf("Commit failed: %v", err)
		return nil, err
	}

	return s.GetMultisigTx(multisigTx.UnsignedTx)
}

func (s *MultisigService) UpdateMultisigTx(multisigTx *model.MultisigTx) (bool, error) {

	if multisig, _ := s.GetMultisigTx(multisigTx.UnsignedTx); multisig == nil {
		return false, errors.New("no pending multisig tx found")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE multisig_tx SET transaction_id = ? WHERE unsigned_tx = ? AND transaction_id IS NULL")
	if err != nil {
		return false, err
	}
	_, err = stmt.Exec(multisigTx.TransactionId, multisigTx.UnsignedTx)
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

func (s *MultisigService) GetAllMultisigTx() (*[]model.MultisigTx, error) {
	return s.doGetMultisigTx("", "")
}

func (s *MultisigService) GetAllMultisigTxForAlias(alias string) (*[]model.MultisigTx, error) {
	tx, err := s.doGetMultisigTx(alias, "")

	if err != nil {
		return nil, err
	}
	if len(*tx) <= 0 {
		return &[]model.MultisigTx{}, nil
	}
	return tx, nil
}

func (s *MultisigService) GetMultisigTx(txId string) (*model.MultisigTx, error) {
	tx, err := s.doGetMultisigTx("", txId)

	if err != nil {
		return nil, err
	}
	if len(*tx) <= 0 {
		return nil, nil
	}

	return &(*tx)[0], nil
}

func (s *MultisigService) AddMultisigTxSigner(txId string, signer *model.MultisigTxSigner) (*model.MultisigTx, error) {
	multisigTx, err := s.GetMultisigTx(txId)
	if err != nil {
		return nil, err
	}

	// check if signer signature or address is empty
	if len(signer.Address) == 0 || len(signer.Signature) == 0 {
		return nil, errors.New("signer address or signature is empty")
	}

	if !s.isOwner(multisigTx, signer.Address) {
		return nil, errors.New("signer is not owner")
	}

	// check if signer count is more than threshold
	signers := multisigTx.Signers
	if len(signers) >= multisigTx.Threshold {
		return nil, errors.New("signer count is more than threshold")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare("UPDATE multisig_tx_owners SET signature = ?, is_signer = ? WHERE multisig_tx_id = ? AND address = ?")
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(signer.Signature, true, multisigTx.Id, signer.Address)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Execute statement failed: %v, unable to rollback: %v", err, rollbackErr)
		}
		log.Print(err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Commit failed: %v", err)
		return nil, err
	}

	return s.GetMultisigTx(txId)
}

func (s *MultisigService) doGetMultisigTx(alias string, txId string) (*[]model.MultisigTx, error) {
	var err error

	query := "SELECT tx.id, " +
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
		"WHERE (tx.alias=? OR ?='') AND (tx.unsigned_tx=? OR ?='') AND tx.transaction_id IS NULL " +
		"ORDER BY tx.created_at ASC"

	rows, err := s.db.Query(query, alias, alias, txId, txId)
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
	var multiSigTx = make(map[int]model.MultisigTx)

	for rows.Next() {
		var (
			txId              int
			txAlias           string
			txThreshold       int
			txTransactionId   sql.NullString
			txUnsignedTx      string
			ownerMultisigTxId sql.NullInt64
			ownerId           sql.NullInt64
			ownerAddress      sql.NullString
			ownerSignature    sql.NullString
			ownerIsSigner     sql.NullBool
		)

		err := rows.Scan(&txId, &txAlias, &txThreshold, &txTransactionId, &txUnsignedTx, &ownerMultisigTxId, &ownerId, &ownerAddress, &ownerSignature, &ownerIsSigner)
		if err != nil {
			log.Fatal(err)
		}
		var tx model.MultisigTx
		if _, ok := multiSigTx[txId]; ok {
			tx = multiSigTx[txId]
		} else {
			tx = model.MultisigTx{
				Id:            int64(txId),
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
		signers := tx.Signers
		if signers == nil {
			signers = []model.MultisigTxSigner{}
		}

		// add owner
		owner := model.MultisigTxOwner{
			Id:           ownerId.Int64,
			MultisigTxId: ownerMultisigTxId.Int64,
			Address:      ownerAddress.String,
		}
		owners = append(owners, owner)

		// add signer
		if ownerIsSigner.Valid && ownerIsSigner.Bool {
			signer := model.MultisigTxSigner{
				MultisigTxOwner: owner,
				Signature:       ownerSignature.String,
			}
			signers = append(signers, signer)
		}

		tx.Owners = owners
		tx.Signers = signers

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

func (s *MultisigService) isOwner(multisigTx *model.MultisigTx, signerAddress string) bool {
	for _, owner := range multisigTx.Owners {
		if owner.Address == signerAddress {
			return true
		}
	}
	return false
}
