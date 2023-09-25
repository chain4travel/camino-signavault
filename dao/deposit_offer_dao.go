package dao

import (
	"database/sql"
	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/model"
	"log"
)

var _ DepositOfferDao = (*depositOfferDao)(nil)

type DepositOfferDao interface {
	AddSignature(depositOfferID, address, signature string) error
	GetSignatures(address string) (*[]model.DepositOfferSig, error)
}
type depositOfferDao struct {
	db *db.Db
}

func NewDepositOfferDao(db *db.Db) DepositOfferDao {
	return &depositOfferDao{
		db: db,
	}
}
func (d *depositOfferDao) AddSignature(depositOfferID, address, signature string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO deposit_offer_sigs (deposit_offer_id, address, signature) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(depositOfferID, address, signature)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Execute statement failed: %v, unable to rollback: %v", err, rollbackErr)
		}
		log.Print(err)
		return err
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
		"sigs.address, " +
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
			depositOfferId string
			signature      string
		)

		err = rows.Scan(&depositOfferId, &address, &signature)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, model.DepositOfferSig{
			DepositOfferID: depositOfferId,
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
