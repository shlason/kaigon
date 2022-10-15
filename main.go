package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shlason/kaigon/configs"
	"github.com/shlason/kaigon/middlewares"
	_ "github.com/shlason/kaigon/models"
	"github.com/shlason/kaigon/routes"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"
)

func main() {
	var g errgroup.Group

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

	//r.Run()

	g.Go(func() error {
		return http.ListenAndServe(":http", http.RedirectHandler(fmt.Sprintf("https://%s", configs.Server.Host), http.StatusSeeOther))
	})
	g.Go(func() error {
		return http.Serve(autocert.NewListener(configs.Server.Host), r)
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
