package controller

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func Proxy(ctx *gin.Context) {
	remote, err := url.Parse("http://google.com")
	if err != nil {
		panic(err)
	}

	internalAuthToken, ok := ctx.Get("INTERNAL_AUTH_TOKEN")
	if !ok {
		panic("Internal server error")
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = ctx.Request.Header
		req.Header.Set("INTERNAL_AUTH_TOKEN", internalAuthToken.(string))

		req.Host = remote.Host

		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = remote.Path

		// NOTE: To log the request body, maybe use github.com/uber-go/zap
		value, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf(string(value))
	}

	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
