package controller

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func Proxy(serviceAddr string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		remote, err := url.Parse(serviceAddr + path)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		internalAuthToken, ok := ctx.Get("Internal_auth_token")
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = remote.Path

			req.Header = ctx.Request.Header
			req.Header.Set("Internal_auth_token", internalAuthToken.(string))

			req.Body = ctx.Request.Body

			// NOTE: To log the request body, maybe use github.com/uber-go/zap
			// value, err := io.ReadAll(req.Body)
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }
			// fmt.Printf(string(value))
		}

		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
