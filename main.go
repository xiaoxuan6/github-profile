package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/xiaoxuan6/github-profile/github"
	"github.com/xiaoxuan6/github-profile/routers"
)

func main() {
	_ = godotenv.Load()
	github.Init()
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")
	routers.RegisterRouter(r)
	_ = r.Run(":11080")
}
