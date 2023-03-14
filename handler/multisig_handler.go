package handler

import (
	"fmt"
	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/dto"
	"github.com/chain4travel/camino-signavault/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MultisigHandler struct {
	MultisigSvc *service.MultisigService
}

func NewMultisigHandler() *MultisigHandler {
	return &MultisigHandler{
		MultisigSvc: service.NewMultisigService(*db.GetInstance()),
	}
}

func (h *MultisigHandler) CreateMultisigTx(ctx *gin.Context) {
	var args *dto.MultisigTxArgs
	err := ctx.BindJSON(&args)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error parsing multisig transaction from JSON",
			"error":   err.Error(),
		})
	}

	response, err := h.MultisigSvc.CreateMultisigTx(args)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error creating new multisig transaction",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(201, *response)
}

func (h *MultisigHandler) GetAllMultisigTx(ctx *gin.Context) {
	multisigTx, err := h.MultisigSvc.GetAllMultisigTx()
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error getting all multisig transactions",
			"error":   err.Error(),
		})
	}
	ctx.JSON(200, multisigTx)
}

func (h *MultisigHandler) GetMultisigTx(ctx *gin.Context) {
	id := h.parseIdParam(ctx.Param("txId"), ctx)

	multisigTx, err := h.MultisigSvc.GetMultisigTx(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("Error getting multisig transaction with id %s", id),
			"error":   err.Error(),
		})
	}

	if multisigTx == nil {
		ctx.JSON(404, gin.H{
			"message": fmt.Sprintf("No pending multisig transaction found with id %s", id),
			"error":   "not found",
		})
		return
	}
	ctx.JSON(200, multisigTx)
}

func (h *MultisigHandler) GetAllMultisigTxForAlias(ctx *gin.Context) {
	alias := ctx.Param("alias")
	multisigTx, err := h.MultisigSvc.GetAllMultisigTxForAlias(alias)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("Error getting all multisig transactions for alias %s", alias),
			"error":   err.Error(),
		})
	}
	if multisigTx == nil {
		ctx.JSON(404, gin.H{
			"message": fmt.Sprintf("Multisig transactions not found for alias %s", alias),
			"error":   "not found",
		})
		return
	}
	ctx.JSON(200, multisigTx)
}

func (h *MultisigHandler) AddMultisigTxSigner(ctx *gin.Context) {
	var err error
	txId := h.parseIdParam(ctx.Param("txId"), ctx)

	var signer *dto.SignTxArgs
	err = ctx.BindJSON(&signer)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error parsing signer from JSON",
			"error":   err.Error(),
		})
		return
	}

	_, err = h.MultisigSvc.AddMultisigTxSigner(txId, signer)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("Error adding signer to multisig transaction with id %s", txId),
			"error":   err.Error(),
		})
		return
	}
	multisigAlias, _ := h.MultisigSvc.GetMultisigTx(txId)
	ctx.JSON(200, multisigAlias)
}

func (h *MultisigHandler) CompleteMultisigTx(ctx *gin.Context) {
	txId := h.parseIdParam(ctx.Param("txId"), ctx)
	var completeTxArgs *dto.CompleteTxArgs
	err := ctx.BindJSON(&completeTxArgs)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error parsing JSON for completing multisig transaction",
			"error":   err.Error(),
		})
	}

	_, err = h.MultisigSvc.UpdateMultisigTx(txId, completeTxArgs)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error completing multisig transaction",
			"error":   err.Error(),
		})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *MultisigHandler) parseIdParam(idParam string, ctx *gin.Context) int64 {
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error parsing id " + ctx.Param("id") + " to integer",
			"error":   err.Error(),
		})
		return 0
	}
	return int64(id)
}
