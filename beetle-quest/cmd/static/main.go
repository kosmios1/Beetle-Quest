package main

import (
	"embed"
	_ "embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"beetle-quest/pkg/httpserver"
	packageMiddleware "beetle-quest/pkg/middleware/secure"
)

//go:embed static
var staticFiles embed.FS

func main() {
	httpserver.GenOwnCertAndKey("static-service")
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(packageMiddleware.NewSecureMiddleware())

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	r.StaticFS("/static", http.FS(staticFS))

	server := httpserver.SetupHTPPSServer(r)
	httpserver.ListenAndServe(server)
}
