package dto

type AddSignatureArgs struct {
	DepositOfferID string `json:"depositOfferID"  binding:"required"`
	Address        string `json:"address"  binding:"required"`
	Signature      string `json:"signature"  binding:"required"`
	Timestamp      int64  `json:"timestamp"` // used for querying deposit offers. optional: if not provided, current time is used
}
