package main

import (
	"github.com/chain4travel/camino-signavault/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	//gin.SetMode(gin.DebugMode)
	router := gin.Default()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		return
	}
	api := router.Group("/v1")

	handler := handlers.NewMultisigHandler()

	api.GET("", handler.GetAllMultisigTx)
	api.GET("/:id", handler.GetMultisigTx)
	api.POST("", handler.CreateMultisigTx)
	api.PUT("/:id", handler.UpdateMultisigTx)

	log.Println("Listing for requests at http://localhost:9000/v1")
	// listen and serve on 0.0.0.0:9000 (for windows "localhost:9000")
	router.Run(":9000")

}
