package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/middlewares"
	_ "github.com/shlason/kaigon/models"
	"github.com/shlason/kaigon/routes"
)

func main() {
	r := gin.Default()

	// Common Middlewares
	r.Use(cors.Default())
	r.Use(gin.Recovery())

	// Public API
	public := r.Group("/api")

	// Private API
	private := r.Group("/api")
	private.Use(middlewares.JWT)

	routes.RegisteAccountRoutes(public, private)
	routes.RegisteAuthRoutes(public)

	r.Run()
}
