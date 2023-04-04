/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package handler

import (
	"fmt"
	"net/http"

	"github.com/chain4travel/camino-signavault/dao"
	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/dto"
	"github.com/chain4travel/camino-signavault/service"
	"github.com/chain4travel/camino-signavault/util"
	"github.com/gin-gonic/gin"
)

type MultisigHandler struct {
	multisigService service.MultisigService
}

func NewMultisigHandler() *MultisigHandler {
	config := util.GetInstance()
	return &MultisigHandler{
		multisigService: service.NewMultisigService(config, dao.NewMultisigTxDao(db.GetInstance()), service.NewNodeService(config)),
	}
}

// CreateMultisigTx godoc
// @Summary Create a new multisig transaction
// @Tags Multisig
// @Accept  json
// @Produce  json
// @Param multisigTxArgs body dto.MultisigTxArgs true "The input parameters for the multisig transaction"
// @Success 201 {object} model.MultisigTx
// @Failure 400 {object} dto.SignavaultError
// @ID CreateMultisigTx
// @Router /multisig [post]
func (h *MultisigHandler) CreateMultisigTx(ctx *gin.Context) {
	var args *dto.MultisigTxArgs
	err := ctx.BindJSON(&args)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			&dto.SignavaultError{
				Message: "Error parsing multisig transaction from JSON",
				Error:   err.Error(),
			})
		return
	}

	response, err := h.multisigService.CreateMultisigTx(args)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			&dto.SignavaultError{
				Message: "Error parsing multisig transaction",
				Error:   err.Error(),
			})
		return
	}
	ctx.JSON(http.StatusCreated, *response)
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
// @Router /multisig/{alias} [get]
func (h *MultisigHandler) GetAllMultisigTxForAlias(ctx *gin.Context) {
	alias := ctx.Param("alias")
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

	multisigTx, err := h.multisigService.GetAllMultisigTxForAlias(alias, timestamp, signature)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			&dto.SignavaultError{
				Message: fmt.Sprintf("Error getting all multisig transactions for alias %s", alias),
				Error:   err.Error(),
			})
		return
	}
	ctx.JSON(http.StatusOK, multisigTx)
}

// SignMultisigTx godoc
// @Summary Signs a multisig transaction
// @Tags Multisig
// @Accept json
// @Produce  json
// @Param id path string true "Multisig transaction ID"
// @Param signTxArgs body dto.SignTxArgs true "Signer details"
// @Success 200 {object} model.MultisigTx
// @Failure 400 {object} dto.SignavaultError
// @ID SignMultisigTx
// @Router /multisig/{id} [put]
func (h *MultisigHandler) SignMultisigTx(ctx *gin.Context) {
	var err error
	id := ctx.Param("id")

	var signer *dto.SignTxArgs
	err = ctx.BindJSON(&signer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			&dto.SignavaultError{
				Message: "Error parsing signer from JSON",
				Error:   err.Error(),
			})
		return
	}

	_, err = h.multisigService.SignMultisigTx(id, signer)
	if err != nil {
		code := http.StatusBadRequest
		if err == service.ErrTxNotExists {
			code = http.StatusNotFound
		}
		ctx.JSON(code,
			&dto.SignavaultError{
				Message: fmt.Sprintf("Error adding signer to multisig transaction with id %s", id),
				Error:   err.Error(),
			})
		return
	}
	multisigAlias, _ := h.multisigService.GetMultisigTx(id)
	ctx.JSON(http.StatusOK, multisigAlias)
}

// IssueMultisigTx issues a new multisig transaction with the given parameters.
// @Summary Issue a new multisig transaction
// @Tags Multisig
// @Accept json
// @Produce json
// @Param issueTxArgs body dto.IssueTxArgs true "IssueTxArgs object that contains the parameters for the multisig transaction to be issued"
// @Success 200 {object} dto.IssueTxResponse
// @Failure 400 {object} dto.SignavaultError
// @ID IssueMultisigTx
// @Router /multisig/issue [post]
func (h *MultisigHandler) IssueMultisigTx(ctx *gin.Context) {
	var issueTxArgs *dto.IssueTxArgs
	err := ctx.BindJSON(&issueTxArgs)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			&dto.SignavaultError{
				Message: "Error parsing JSON for issuing multisig transaction",
				Error:   err.Error(),
			})
		return
	}

	txID, err := h.multisigService.IssueMultisigTx(issueTxArgs)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			&dto.SignavaultError{
				Message: "Error issuing multisig transaction",
				Error:   err.Error(),
			})
		return
	}
	ctx.JSON(http.StatusOK, &dto.IssueTxResponse{TxID: txID.String()})
}

func (h *MultisigHandler) throwMissingQueryParamError(ctx *gin.Context, param string) {
	ctx.JSON(http.StatusBadRequest,
		&dto.SignavaultError{
			Message: fmt.Sprintf("Missing query parameter '%s'", param),
			Error:   "missing query parameter",
		})
}
