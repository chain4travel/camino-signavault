package handler

import (
	"fmt"
	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/model"
	"github.com/chain4travel/camino-signavault/service"
	"github.com/gin-gonic/gin"
)

type MultisigHandler struct {
	MultisigSvc *service.MultisigService
}

func NewMultisigHandler() *MultisigHandler {
	return &MultisigHandler{
		MultisigSvc: service.NewMultisigService(*db.GetInstance()),
	}
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
			"error":   "Not Found",
		})
		return
	}
	ctx.JSON(200, multisigTx)
}
func (h *MultisigHandler) GetMultisigTx(ctx *gin.Context) {
	id := ctx.Param("txId")

	multisigTx, err := h.MultisigSvc.GetMultisigTx(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("Error getting multisig transaction with id %s", id),
			"error":   err.Error(),
		})
	}

	if multisigTx == nil {
		ctx.JSON(404, gin.H{
			"message": fmt.Sprintf("Multisig transaction not found for id %s", id),
			"error":   "Not Found",
		})
		return
	}
	ctx.JSON(200, multisigTx)
}

func (h *MultisigHandler) CreateMultisigTx(ctx *gin.Context) {
	var multisigTx *model.MultisigTx
	err := ctx.BindJSON(&multisigTx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error parsing multisig transaction from JSON",
			"error":   err.Error(),
		})
	}

	response, err := h.MultisigSvc.CreateMultisigTx(multisigTx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error inserting multisig transaction in database",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(200, *response)
}

func (h *MultisigHandler) UpdateMultisigTx(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "UpdateMultisigTx",
		"error":   "Not implemented",
	})
}

func (h *MultisigHandler) AddMultisigTxSigner(context *gin.Context) {
	var err error
	txId := context.Param("txId")

	var signer *model.MultisigTxSigner
	err = context.BindJSON(&signer)
	if err != nil {
		context.JSON(400, gin.H{
			"message": "Error parsing signer from JSON",
			"error":   err.Error(),
		})
		return
	}

	_, err = h.MultisigSvc.AddMultisigTxSigner(txId, signer)
	if err != nil {
		context.JSON(400, gin.H{
			"message": fmt.Sprintf("Error adding signer %s to multisig transaction with id %d", signer.Address, txId),
			"error":   err.Error(),
		})
		return
	}
	multisigAlias, _ := h.MultisigSvc.GetMultisigTx(txId)
	context.JSON(200, multisigAlias)
}
