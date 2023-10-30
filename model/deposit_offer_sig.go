package model

type DepositOfferSig struct {
	DepositOfferID string `json:"depositOfferID"`
	Address        string `json:"address"`
	Signature      string `json:"signature"`
}
