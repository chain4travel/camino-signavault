package handler

import (
	"fmt"
	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/model"
	"github.com/chain4travel/camino-signavault/service"
	"github.com/gin-gonic/gin"
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
	alias := ctx.Param("alias")
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error parsing id " + ctx.Param("id") + " to integer",
			"error":   err.Error(),
		})
		return
	}

	multisigTx, err := h.MultisigSvc.GetMultisigTx(alias, id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("Error getting multisig transaction with id %d for alias %s", id, alias),
			"error":   err.Error(),
		})
	}

	if multisigTx == nil {
		ctx.JSON(404, gin.H{
			"message": fmt.Sprintf("Multisig transaction not found for alias %s and id %d ", alias, id),
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

	id, err := h.MultisigSvc.CreateMultisigTx(multisigTx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error inserting multisig transaction in database",
			"error":   err.Error(),
		})
		return
	}
	multisigTx.Id = id
	ctx.JSON(200, *multisigTx)
}

func (h *MultisigHandler) UpdateMultisigTx(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "UpdateMultisigTx",
		"error":   "Not implemented",
	})
}

func (h *MultisigHandler) AddMultisigTxSigner(context *gin.Context) {
	alias := context.Param("alias")
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(400, gin.H{
			"message": "Error parsing id " + context.Param("id") + " to integer",
			"error":   err.Error(),
		})
		return
	}

	var signer *model.MultisigTxSigner
	err = context.BindJSON(&signer)
	if err != nil {
		context.JSON(400, gin.H{
			"message": "Error parsing signer from JSON",
			"error":   err.Error(),
		})
		return
	}

	_, err = h.MultisigSvc.AddMultisigTxSigner(id, signer)
	if err != nil {
		context.JSON(400, gin.H{
			"message": fmt.Sprintf("Error adding signer %s to multisig transaction with id %d for alias %s", signer.Address, id, alias),
			"error":   err.Error(),
		})
		return
	}
	multisigAlias, _ := h.MultisigSvc.GetMultisigTx(alias, id)
	context.JSON(200, multisigAlias)
}
