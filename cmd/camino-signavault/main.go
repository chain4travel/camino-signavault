package main

import (
	"github.com/chain4travel/camino-signavault/handler"
	"github.com/chain4travel/camino-signavault/util"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	config := util.GetInstance()
	startRouter(config)
}

func startRouter(cfg *util.Config) {
	//gin.SetMode(gin.DebugMode)
	router := gin.Default()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		return
	}
	api := router.Group("/v1")

	h := handler.NewMultisigHandler()

	api.GET("/multisig", h.GetAllMultisigTx)
	api.GET("/multisig/:alias", h.GetAllMultisigTxForAlias)
	api.GET("/multisig/:alias/:id", h.GetMultisigTx)
	api.POST("/multisig", h.CreateMultisigTx)
	api.POST("/multisig/:alias/:id", h.AddMultisigTxSigner)

	err = router.Run(cfg.ListenerAddress)
	if err != nil {
		log.Fatal(err)
	}
}
