package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/grigoriy-st/YL-Golang/docs"
	"github.com/grigoriy-st/YL-Golang/internal/orchestrator/handler"
	"github.com/grigoriy-st/YL-Golang/internal/orchestrator/middlewares"
	"github.com/grigoriy-st/YL-Golang/pkg/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(router *gin.Engine) *gin.Engine {
	v1 := router.Group("/api/v1")
	v1.Use(middlewares.CORS())
	{
		{
			auth := &handler.Auth{Route: v1.Group("/")}
			auth.Route.POST("/register", auth.Register)
			auth.Route.POST("/login", auth.Login)

			auth.Route.OPTIONS("/register", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})
			auth.Route.OPTIONS("/login", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})
		}

		// task
		{
			task := &handler.Task{Route: v1.Group("/task")}
			task.Route.Use(middlewares.Auth())

			task.Route.GET("", task.Index)
			task.Route.POST("", task.Store)
			task.Route.GET("/:id", task.Show)

			task.Route.OPTIONS("", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})
		}

		{
			agent := &handler.Agent{Route: v1.Group("/agent")}
			agent.Route.GET("", agent.Index)
			agent.Route.GET("/ws", agent.WebSocket)
		}
	}

	if config.Get().Mode == "debug" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DefaultModelsExpandDepth(-1)))
	}

	return router
}
