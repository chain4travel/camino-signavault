package service

import (
	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/model"
	"github.com/gocraft/dbr/v2"
)

type MultisigService struct {
}

func NewMultisigService() *MultisigService {
	return &MultisigService{}
}

func (s *MultisigService) CreateMultisigTx(multisigTx *model.MultisigTx) (int64, error) {
	var err error
	session := db.GetInstance().NewSession(nil)
	tx, err := session.Begin()

	defer tx.RollbackUnlessCommitted()

	exec, err := tx.InsertInto("multisig_tx").
		Pair("alias", multisigTx.Alias).
		Pair("threshold", multisigTx.Threshold).
		Pair("transaction_id", multisigTx.TransactionId).
		Pair("unsigned_tx", multisigTx.UnsignedTx).
		Exec()
	if err != nil {
		return 0, err
	}

	id, _ := exec.LastInsertId()

	for _, signer := range multisigTx.Signers {
		exec, err = tx.InsertInto("multisig_tx_signers").
			Pair("multisig_tx_id", id).
			Pair("address", signer.Address).
			Pair("signature", signer.Signature).
			Exec()
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit()

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *MultisigService) GetMultisigTx(alias string) (*[]model.MultisigTx, error) {
	session := db.GetInstance().NewSession(nil)
	var err error

	_, err = session.Begin()

	multisigTx := &[]model.MultisigTx{}

	_, err = session.
		Select("multisig_tx.id", "alias", "threshold", "transaction_id", "unsigned_tx", "signers.multisig_tx_id", "signers.id", "signers.address", "signers.signature").
		From("multisig_tx").
		LeftJoin(dbr.I("multisig_tx_signers").As("signers"), "multisig_tx.id = signers.multisig_tx_id").
		Where("alias=?", alias).
		Load(multisigTx)

	if err != nil {
		return nil, err
	}

	return multisigTx, nil
}
