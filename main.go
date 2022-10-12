package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/routes"
)

func main() {
	r := gin.Default()

	r.Use(cors.Default())
	apiRoutes := r.Group("/api")

	routes.RegisteAccountRoutes(apiRoutes)
	routes.RegisteAuthRoutes(apiRoutes)

	r.Run()
}
