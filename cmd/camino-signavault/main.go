/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package main

import (
	"log"

	"github.com/chain4travel/camino-signavault/dao"
	"github.com/chain4travel/camino-signavault/db"
	"github.com/chain4travel/camino-signavault/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/chain4travel/camino-signavault/handler"
	"github.com/chain4travel/camino-signavault/util"
)

// Command to generate swagger docs (root dir)
// swag init -g ./cmd/camino-signavault/main.go --exclude ./dependencies/caminogo
// Command to generate the typescript client
// openapi-generator-cli generate -i docs/swagger.json -g typescript-axios -o signavaultjs

// @title Signavault API
// @version 1.0
// @description This is the signavault API.
// @host localhost:8080
// @BasePath /v1
// @schemes http
func main() {
	config := util.GetInstance()
	startRouter(config)
}

func startRouter(cfg *util.Config) {
	// gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Use(cors.Default())
	err := router.SetTrustedProxies(nil)
	if err != nil {
		return
	}
	api := router.Group("/v1")

	multisigService := service.NewMultisigService(cfg, dao.NewMultisigTxDao(db.GetInstance()), service.NewNodeService(cfg))
	h := handler.NewMultisigHandler(multisigService)

	api.POST("/multisig", h.CreateMultisigTx)
	api.POST("/multisig/issue", h.IssueMultisigTx)
	api.POST("/multisig/cancel", h.CancelMultisigTx)
	api.PUT("/multisig/:id", h.SignMultisigTx)
	api.GET("/multisig/:alias", h.GetAllMultisigTxForAlias)

	err = router.Run(cfg.ListenerAddress)
	if err != nil {
		log.Fatal(err)
	}
}
