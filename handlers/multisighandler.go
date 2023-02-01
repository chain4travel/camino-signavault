package handlers

import (
	"github.com/chain4travel/camino-signavault/model"
	"github.com/gin-gonic/gin"
)

type MultisigHandler struct {
}

func NewMultisigHandler() *MultisigHandler {
	return &MultisigHandler{}
}

func (h *MultisigHandler) GetAllMultisigTx(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "GetAllMultisigTx",
	})
}
func (h *MultisigHandler) GetMultisigTx(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "GetMultisigTx",
	})
}

func (h *MultisigHandler) CreateMultisigTx(ctx *gin.Context) {
	var multisigTx *model.MultisigTx
	err := ctx.BindJSON(&multisigTx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error parsing MultisigTx JSON",
		})
	}
	ctx.JSON(200, *multisigTx)
}

func (h *MultisigHandler) UpdateMultisigTx(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "UpdateMultisigTx",
	})
}
