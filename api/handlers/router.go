package router

import "github.com/gin-gonic/gin"

func Router() *gin.Engine {
	router := gin.Default()
	router.StaticFS("/static", http.Dir("./Static"))
	router.Use(gin.Recovery())
	router.Use(cors.AllowAll())

	router.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"message": "API working fine!",
			},
		)
	})

	superGroup := router.Group("/api/v1/") {
		userGroup := superGroup.Group("/user/") {
			userGroup.POST("register", )
		}
	}
}
