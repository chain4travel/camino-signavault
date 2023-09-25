/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package handler

import (
	"fmt"
	"net/http"

	"github.com/chain4travel/camino-signavault/dto"
	"github.com/chain4travel/camino-signavault/service"
	"github.com/gin-gonic/gin"
)

type DepositOfferHandler interface {
	AddSignature(ctx *gin.Context)
	GetSignatures(ctx *gin.Context)
}

type depositOfferHandler struct {
	DepositOfferService service.DepositOfferService
}

func NewDepositOfferHandler(DepositOfferService service.DepositOfferService) *depositOfferHandler {
	return &depositOfferHandler{
		DepositOfferService: DepositOfferService,
	}
}

// AddSignature godoc
// @Summary Adds a signature mapped to a deposit offer id and an address
// @Tags Multisig
// @Accept  json
// @Produce  json
// @Param addSignatureArgs body dto.AddSignatureArgs true "The input parameters for the multisig transaction"
// @Success 201
// @Failure 400 {object} dto.SignavaultError
// @ID AddSignature
// @Router /deposit-offer [post]
func (h *depositOfferHandler) AddSignature(ctx *gin.Context) {
	var args *dto.AddSignatureArgs
	err := ctx.BindJSON(&args)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			&dto.SignavaultError{
				Message: "Error parsing signature args from JSON",
				Error:   err.Error(),
			})
		return
	}

	err = h.DepositOfferService.AddSignature(args)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			&dto.SignavaultError{
				Message: "Error creating multisig transaction",
				Error:   err.Error(),
			})
		return
	}
	ctx.Status(http.StatusCreated)
}

// GetAllMultisigTxForAlias godoc
// @Summary Retrieves all multisig transactions for a given alias
// @Tags Multisig
// @Param alias path string true "Alias of the multisig account"
// @Param signature query string true "Signature for the request"
// @Param timestamp query string true "Timestamp for the request"
// @Produce  json
// @Success 200 {array} model.MultisigTx
// @Failure 400 {object}  dto.SignavaultError
// @ID GetAllMultisigTxForAlias
// @Router /deposit-offer/{address} [get]
func (h *depositOfferHandler) GetSignatures(ctx *gin.Context) {
	address := ctx.Param("address")
	signature, b := ctx.GetQuery("signature")
	if !b {
		h.throwMissingQueryParamError(ctx, "signature")
		return
	}
	timestamp, b := ctx.GetQuery("timestamp")
	if !b {
		h.throwMissingQueryParamError(ctx, "timestamp")
		return
	}

	sigs, err := h.DepositOfferService.GetSignatures(address, timestamp, signature)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			&dto.SignavaultError{
				Message: fmt.Sprintf("Error getting all deposit offer signatures for adadress %s", address),
				Error:   err.Error(),
			})
		return
	}
	ctx.JSON(http.StatusOK, sigs)
}

func (h *depositOfferHandler) throwMissingQueryParamError(ctx *gin.Context, param string) {
	ctx.JSON(http.StatusBadRequest,
		&dto.SignavaultError{
			Message: fmt.Sprintf("Missing query parameter '%s'", param),
			Error:   "missing query parameter",
		})
}
