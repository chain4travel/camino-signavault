package handler

import (
	"github.com/chain4travel/camino-signavault/model"
	"github.com/chain4travel/camino-signavault/service"
	"github.com/gin-gonic/gin"
)

type MultisigHandler struct {
	MultisigSvc *service.MultisigService
}

func NewMultisigHandler() *MultisigHandler {
	return &MultisigHandler{MultisigSvc: service.NewMultisigService()}
}

func (h *MultisigHandler) GetAllMultisigTx(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "GetAllMultisigTx",
	})
}
func (h *MultisigHandler) GetMultisigTx(ctx *gin.Context) {
	alias := ctx.Param("alias")
	multisigTx, err := h.MultisigSvc.GetMultisigTx(alias)

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error getting MultisigTx for alias " + alias,
			"error":   err.Error(),
		})
	}

	ctx.JSON(200, multisigTx)
}

func (h *MultisigHandler) CreateMultisigTx(ctx *gin.Context) {
	var multisigTx *model.MultisigTx
	err := ctx.BindJSON(&multisigTx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error parsing MultisigTx from JSON",
			"error":   err.Error(),
		})
	}

	id, err := h.MultisigSvc.CreateMultisigTx(multisigTx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error inserting MultisigTx in database",
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
