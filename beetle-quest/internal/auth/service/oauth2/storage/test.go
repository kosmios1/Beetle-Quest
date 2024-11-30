//go:build beetleQuestTest

package storage

import (
	"log"

	"github.com/go-oauth2/oauth2/v4"
	o2store "github.com/go-oauth2/oauth2/v4/store"
)

func GetTokenStorage() oauth2.TokenStore {
	if store, err := o2store.NewMemoryTokenStore(); err != nil {
		log.Fatalf("[FATAL] Could not setup oauth2 token storage!")
		return nil
	} else {
		return store
	}
}
