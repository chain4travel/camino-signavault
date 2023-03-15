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

	api.POST("/multisig", h.CreateMultisigTx)
	api.POST("/multisig/:txId", h.CompleteMultisigTx)
	api.PUT("/multisig/:txId", h.SignMultisigTx)
	api.GET("/multisig/:alias", h.GetAllMultisigTxForAlias)

	err = router.Run(cfg.ListenerAddress)
	if err != nil {
		log.Fatal(err)
	}
}
