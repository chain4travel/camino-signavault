package handler

import (
	"fmt"
	"github.com/chain4travel/camino-signavault/dao"
	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/dto"
	"github.com/chain4travel/camino-signavault/service"
	"github.com/chain4travel/camino-signavault/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

func (h *MultisigHandler) CreateMultisigTx(ctx *gin.Context) {
	var args *dto.MultisigTxArgs
	err := ctx.BindJSON(&args)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error parsing multisig transaction from JSON",
			"error":   err.Error(),
		})
		return
	}

	response, err := h.multisigService.CreateMultisigTx(args)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error creating new multisig transaction",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(201, *response)
}

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
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("Error getting all multisig transactions for alias %s", alias),
			"error":   err.Error(),
		})
		return
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

func (h *MultisigHandler) SignMultisigTx(ctx *gin.Context) {
	var err error
	id := ctx.Param("id")

	var signer *dto.SignTxArgs
	err = ctx.BindJSON(&signer)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error parsing signer from JSON",
			"error":   err.Error(),
		})
		return
	}

	_, err = h.multisigService.SignMultisigTx(id, signer)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("Error adding signer to multisig transaction with id %s", id),
			"error":   err.Error(),
		})
		return
	}
	multisigAlias, _ := h.multisigService.GetMultisigTx(id)
	ctx.JSON(200, multisigAlias)
}

func (h *MultisigHandler) CompleteMultisigTx(ctx *gin.Context) {
	id := ctx.Param("id")
	var completeTxArgs *dto.CompleteTxArgs
	err := ctx.BindJSON(&completeTxArgs)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error parsing JSON for completing multisig transaction",
			"error":   err.Error(),
		})
		return
	}

	_, err = h.multisigService.CompleteMultisigTx(id, completeTxArgs)
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

func (h *MultisigHandler) throwMissingQueryParamError(ctx *gin.Context, param string) {
	ctx.JSON(400, gin.H{
		"message": fmt.Sprintf("Missing query parameter '%s'", param),
		"error":   "missing query parameter",
	})
}
