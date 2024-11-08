package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaoxuan6/github-profile/handlers"
)

func RegisterRouter(g *gin.Engine) {
	g.GET("/", handlers.IndexHandler.Index)
	g.POST("/generate", handlers.IndexHandler.Generate)

	g.NoMethod(func(c *gin.Context) {
		c.Redirect(301, "/")
	})
	g.NoRoute(func(c *gin.Context) {
		c.Redirect(301, "/")
	})
}
