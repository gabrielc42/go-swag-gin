package router

import (
	"net/http"

	todoController "github.com/gabrielc42"

	"github.com/gin-gonic/gin"
)

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

	superGroup := router.Group("/api/v1/")
	{
		userGroup := superGroup.Group("/user/")
		{
			userGroup.POST("register", userController.Register)
			userGroup.POST("login", userController.Login)
		}
		todoGroup := superGroup.Group("/todo/")
		{
			todoGroup.Use(middleware.TokenAuthMiddleware())
			{
				todoGroup.GET("getTodos", todoController.GetTodos)                //get TODOs
				todoGroup.POST("create", todoController.CreateTodo)               //create TODO
				todoGroup.GET("getTodo/:todoId", todoController.GetTodo)          //create TODO
				todoGroup.PATCH("updateTodo/:todoId", todoController.UpdateTodo)  //get TODOs
				todoGroup.DELETE("deleteTodo/:todoId", todoController.DeleteTodo) //get TODOs
			}

			todoGroup.GET("/hello", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "TODO",
				})
			})
		}
	}
}
