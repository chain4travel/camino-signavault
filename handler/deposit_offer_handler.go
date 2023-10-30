/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package handler

import (
	"fmt"
	"net/http"
	"strconv"

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
// @Tags DepositOffer
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
				Message: "Error parsing signature args",
				Error:   err.Error(),
			})
		return
	}

	err = h.DepositOfferService.AddSignatures(args)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			&dto.SignavaultError{
				Message: "Error inserting signature",
				Error:   err.Error(),
			})
		return
	}
	ctx.Status(http.StatusCreated)
}

// GetSignatures godoc
// @Summary Retrieves all signatures for an address only for authorized calls.
// @Tags DepositOffer
// @Param address path string true "Address for which to retrieve all signatures"
// @Param signature query string true "Signature for the request"
// @Param timestamp query string true "Timestamp for the request"
// @Param multisig query string true "true if the address is a multisig address, false otherwise"
// @Produce  json
// @Success 200 {array} model.DepositOfferSig
// @Failure 400 {object}  dto.SignavaultError
// @ID GetSignatures
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
	multisigParam, b := ctx.GetQuery("multisig")
	if !b {
		h.throwMissingQueryParamError(ctx, "multisig")
		return
	}
	multisig, err := strconv.ParseBool(multisigParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &dto.SignavaultError{
			Message: fmt.Sprintf("invalid query parameter: multisig='%s'", multisigParam),
			Error:   "invalid query parameter",
		})
		return
	}

	sigs, err := h.DepositOfferService.GetSignatures(address, timestamp, signature, multisig)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			&dto.SignavaultError{
				Message: fmt.Sprintf("Error getting all deposit offer signatures for address %s", address),
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
