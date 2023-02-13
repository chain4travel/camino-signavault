package main

import (
	"github.com/chain4travel/camino-signavault/handler"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	readConfig()
	startRouter()
}

func readConfig() {

}

func startRouter() {
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

	log.Println("Listening for requests at http://localhost:9000/v1/multisig")
	// listen and serve on 0.0.0.0:9000 (for windows "localhost:9000")
	err = router.Run(":9000")
	if err != nil {
		log.Fatal(err)
	}
}
