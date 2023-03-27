/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package main

import (
	"log"

	"github.com/chain4travel/camino-signavault/handler"
	"github.com/chain4travel/camino-signavault/util"
	"github.com/gin-gonic/gin"
)

func main() {
	config := util.GetInstance()
	startRouter(config)
}

func startRouter(cfg *util.Config) {
	// gin.SetMode(gin.DebugMode)
	router := gin.Default()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		return
	}
	api := router.Group("/v1")

	h := handler.NewMultisigHandler()

	api.POST("/multisig", h.CreateMultisigTx)
	api.POST("/multisig/issue", h.IssueMultisigTx)
	api.PUT("/multisig/:id", h.SignMultisigTx)
	api.GET("/multisig/:alias", h.GetAllMultisigTxForAlias)

	err = router.Run(cfg.ListenerAddress)
	if err != nil {
		log.Fatal(err)
	}
}
