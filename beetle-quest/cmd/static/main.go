package main

import (
	"embed"
	_ "embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed static
var staticFiles embed.FS

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())
	// TODO: Uncomment this when having a valid SSL certificate
	// r.Use(secure.New(secure.Config{
	// 	SSLRedirect:           true,
	// 	IsDevelopment:         false,
	// 	STSSeconds:            315360000,
	// 	STSIncludeSubdomains:  true,
	// 	FrameDeny:             true,
	// 	ContentTypeNosniff:    true,
	// 	BrowserXssFilter:      true,
	// 	ContentSecurityPolicy: "default-src 'self'",
	// 	IENoOpen:              true,
	// 	SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
	// 	AllowedHosts:          []string{},
	// }))

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	r.Any("/", func(c *gin.Context) {
		c.FileFromFS("static/index.html", http.FS(staticFS))
	})

	r.StaticFS("/static", http.FS(staticFS))

	// imagesFS, err := fs.Sub(imageFiles, "images")
	// if err != nil {
	// 	log.Fatal("Failed to create sub filesystem for images: ", err)
	// }
	// r.StaticFS("/api/v1/static/images", http.FS(imagesFS))

	r.Run()
}
