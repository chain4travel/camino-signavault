package dao

import (
	"database/sql"
	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/model"
	"log"
)

var _ DepositOfferDao = (*depositOfferDao)(nil)

type DepositOfferDao interface {
	AddSignatures(depositOfferID string, addresses, signatures []string) error
	GetSignatures(address string) (*[]model.DepositOfferSig, error)
}
type depositOfferDao struct {
	db             *db.Db
	preparedInsert *sql.Stmt
}

func NewDepositOfferDao(db *db.Db) DepositOfferDao {
	stmt, err := db.Prepare("INSERT INTO deposit_offer_sigs (deposit_offer_id, address, signature) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	return &depositOfferDao{
		db:             db,
		preparedInsert: stmt,
	}

}
func (d *depositOfferDao) AddSignatures(depositOfferID string, addresses []string, signatures []string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Execute statement failed: %v, unable to rollback: %v", err, rollbackErr)
			}
			log.Print(err)
		}
	}()

	for i, address := range addresses {
		_, err = d.preparedInsert.Exec(depositOfferID, address, signatures[i])
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Commit failed: %v", err)
		return err
	}

	return nil
}

func (d *depositOfferDao) GetSignatures(address string) (*[]model.DepositOfferSig, error) {
	var err error

	var rows *sql.Rows
	query := "SELECT sigs.deposit_offer_id, " +
		"sigs.signature " +
		"FROM deposit_offer_sigs AS sigs " +
		"WHERE sigs.address=?"
	rows, err = d.db.Query(query, address)

	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Print(err)
		}
	}(rows)

	var result []model.DepositOfferSig

	for rows.Next() {
		var (
			depositOfferID string
			signature      string
		)

		err = rows.Scan(&depositOfferID, &signature)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, model.DepositOfferSig{
			DepositOfferID: depositOfferID,
			Address:        address,
			Signature:      signature,
		})

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	if result == nil {
		return nil, nil
	}
	return &result, nil

}
