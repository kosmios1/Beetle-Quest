package repository

import (
	"beetle-quest/pkg/models"
	"sync"
)

type SessionRepo struct {
	mux    sync.RWMutex
	tokens map[string]string
}

func NewSessionRepo() (*SessionRepo, error) {
	return &SessionRepo{
		mux:    sync.RWMutex{},
		tokens: make(map[string]string),
	}, nil
}

func (r *SessionRepo) CreateSession(token string) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, ok := r.tokens[token]; !ok {
		r.tokens[token] = token
		return nil
	}
	return models.ErrInternalServerError
}

func (r *SessionRepo) RevokeToken(token string) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, ok := r.tokens[token]; ok {
		delete(r.tokens, token)
		return nil
	}
	return models.ErrInternalServerError
}

func (r *SessionRepo) FindToken(token string) (string, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	if tok, ok := r.tokens[token]; ok {
		return tok, nil
	}
	return "", models.ErrInternalServerError
}
