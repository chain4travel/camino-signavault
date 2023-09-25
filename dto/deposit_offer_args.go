package dto

type AddSignatureArgs struct {
	DepositOfferID string `json:"depositOfferID"`
	Address        string `json:"address"`
	Signature      string `json:"signature"`
}
