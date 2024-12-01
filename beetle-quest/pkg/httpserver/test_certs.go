//go:build beetleQuestTest

package httpserver

import (
	"log"
	"net/http"
)

func GenOwnCertAndKey(serviceName string) {
	log.Println("[INFO] In mock implementation own cert and key are not generated")
}

func SetupHTPPSServer(h http.Handler) *http.Server {
	server := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}
	return server
}

func SetupHTTPSClient() *http.Client {
	client := &http.Client{}
	return client
}

func ListenAndServe(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
