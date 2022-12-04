package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/configs"
	"github.com/shlason/kaigon/docs"
	"github.com/shlason/kaigon/middlewares"
	_ "github.com/shlason/kaigon/models"
	"github.com/shlason/kaigon/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"
)

// @title       Kaigon API
// @version     1.0
// @description This is a forum server.

// @contact.name  API Support
// @contact.url   https://github.com/shlason/kaigon
// @contact.email nocvi111@gmail.com

// @license.name MIT
// @license.url  https://github.com/shlason/kaigon/blob/main/LICENSE

// @host     kaigon.sidesideeffect.io
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
func main() {
	var g errgroup.Group

	docs.SwaggerInfo.Schemes = []string{"https"}

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Common Middlewares
	if os.Getenv("CODE_RUN_ENV") == "prod" {
		r.Use(middlewares.CORS())
	}
	r.Use(gin.Recovery())

	// Public API
	public := r.Group("/api")

	// Private API
	private := r.Group("/api")
	private.Use(middlewares.JWT)

	routes.RegisteAccountRoutes(public, private)
	routes.RegisteAuthRoutes(public, private)
	routes.RegisteForumAndPostRoutes(public, private)
	routes.RegisteTopicRoutes(public, private)
	routes.RegisteImageRoutes(private)
	// TODO: 記得改回 private route
	routes.RegisteChatRoutes(public)
	routes.RegisteSearchRoutes(public)
	routes.RegisteDevelopUtilsRoutes(public)

	if os.Getenv("CODE_RUN_ENV") == "prod" {
		g.Go(func() error {
			return http.ListenAndServe(":http", http.RedirectHandler(fmt.Sprintf("https://%s", configs.Server.Host), http.StatusSeeOther))
		})
		g.Go(func() error {
			return http.Serve(autocert.NewListener(configs.Server.Host), r)
		})

		if err := g.Wait(); err != nil {
			log.Fatal(err)
		}
	} else {
		r.Run()
	}
}
