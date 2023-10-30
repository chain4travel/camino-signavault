package dto

type AddSignatureArgs struct {
	DepositOfferID string   `json:"depositOfferID"  binding:"required"`
	Addresses      []string `json:"addresses"  binding:"required"`
	Signatures     []string `json:"signatures"  binding:"required"`
	Timestamp      int64    `json:"timestamp"` // used for querying deposit offers. optional: if not provided, current time is used
}
