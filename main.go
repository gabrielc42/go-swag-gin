package main

import (
	"fmt"
	"os"

	router "github.com/gabrielc42/api/Handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	r := gin.Default()
	r = router.Router()

	godotenv.Load(".env")
	port := os.Getenv("PORT")
	fmt.Println("this is a port", port)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	fmt.Print("Listening on port")
	r.Run(":" + port)
}
