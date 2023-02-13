package service

import (
	"database/sql"
	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/model"
	"log"
)

type MultisigService struct {
}

func NewMultisigService() *MultisigService {
	return &MultisigService{}
}

func (s *MultisigService) CreateMultisigTx(multisigTx *model.MultisigTx) (int64, error) {
	var err error

	tx, err := db.GetInstance().Begin()
	if err != nil {
		return 0, err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Print(err)
		}
	}(tx)

	stmt, err := tx.Prepare("INSERT INTO multisig_tx (alias, threshold, transaction_id, unsigned_tx) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(multisigTx.Alias, multisigTx.Threshold, multisigTx.TransactionId, multisigTx.UnsignedTx)
	if err != nil {
		return 0, err
	}
	txId, _ := res.LastInsertId()

	for _, signer := range multisigTx.Signers {
		stmt, err := tx.Prepare("INSERT INTO multisig_tx_signers (multisig_tx_id, address, signature) VALUES (?, ?, ?)")
		if err != nil {
			return 0, err
		}
		_, err = stmt.Exec(txId, signer.Address, signer.Signature)
		if err != nil {
			return 0, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return txId, nil
}

func (s *MultisigService) GetAllMultisigTx() (*[]model.MultisigTx, error) {
	return s.doGetMultisigTx("", -1)
}

func (s *MultisigService) GetAllMultisigTxForAlias(alias string) (*[]model.MultisigTx, error) {
	tx, err := s.doGetMultisigTx(alias, -1)

	if err != nil {
		return nil, err
	}
	if len(*tx) <= 0 {
		return nil, nil
	}
	return tx, nil
}

func (s *MultisigService) GetMultisigTx(alias string, id int) (*model.MultisigTx, error) {
	tx, err := s.doGetMultisigTx(alias, id)

	if err != nil {
		return nil, err
	}
	if len(*tx) <= 0 {
		return nil, nil
	}

	return &(*tx)[0], nil
}

func (s *MultisigService) doGetMultisigTx(alias string, id int) (*[]model.MultisigTx, error) {
	var err error

	query := "SELECT tx.id, " +
		"tx.alias, " +
		"tx.threshold, " +
		"tx.transaction_id, " +
		"tx.unsigned_tx, " +
		"signers.multisig_tx_id, " +
		"signers.id, " +
		"signers.address, " +
		"signers.signature " +
		"FROM multisig_tx AS tx " +
		"LEFT JOIN multisig_tx_signers AS signers ON tx.id = signers.multisig_tx_id " +
		"WHERE (tx.alias=? OR ?='') AND (tx.id=? OR ?=-1) " +
		"ORDER BY tx.created_at ASC"

	rows, err := db.GetInstance().Query(query, alias, alias, id, id)
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
			txId               int
			txAlias            string
			txThreshold        int
			txTransactionId    string
			txUnsignedTx       string
			signerMultisigTxId int
			signerId           int
			signerAddress      string
			signerSignature    string
		)

		err := rows.Scan(&txId, &txAlias, &txThreshold, &txTransactionId, &txUnsignedTx, &signerMultisigTxId, &signerId, &signerAddress, &signerSignature)
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
				TransactionId: txTransactionId,
				UnsignedTx:    txUnsignedTx,
			}

		}
		signers := tx.Signers
		if (signers) == nil {
			signers = []model.MultisigTxSigner{}
		}
		signers = append(signers, model.MultisigTxSigner{
			Id:           int64(signerId),
			MultisigTxId: int64(signerMultisigTxId),
			Address:      signerAddress,
			Signature:    signerSignature,
		})
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

func (s *MultisigService) AddMultisigTxSigner(id int, signer *model.MultisigTxSigner) (int64, error) {
	tx, err := db.GetInstance().Begin()
	if err != nil {
		return 0, err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Print(err)
		}
	}(tx)

	stmt, err := tx.Prepare("INSERT INTO multisig_tx_signers (multisig_tx_id, address, signature) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(id, signer.Address, signer.Signature)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}
