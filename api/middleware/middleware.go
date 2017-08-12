// Package middleware contains the HTTP middleware used in the api as
// well as utility functions for interacting with them
package middleware

import (
	"errors"
	"net/http"

	"github.com/praelatus/backend/models"
)

type SessionCache interface {
	Get(id string) (models.Session, error)
	Set(id string, user models.Session) error
	Remove(id string) error
}

type MemCache struct {
	Cache map[string]models.Session
}

func (m *MemCache) Get(id string) (models.Session, error) {
	if val, ok := m.Cache[id]; ok {
		return val, nil
	}

	return models.Session{}, errors.New("No session found")
}

func (m *MemCache) Set(id string, sess models.Session) error {
	m.Cache[id] = sess
	return nil
}

func (m *MemCache) Remove(id string) error {
	delete(m.Cache, id)
	return nil
}

func NewMemCache() *MemCache {
	m := MemCache{}
	m.Cache = make(map[string]models.Session)
	return &m
}

// Cache is the global SessionCache
var Cache SessionCache = NewMemCache()

// LoadMw will wrap the given http.Handler in the DefaultMiddleware
func LoadMw(handler http.Handler) http.Handler {
	h := handler

	for _, m := range DefaultMiddleware {
		h = m(h)
	}

	return h
}

// DefaultMiddleware is the default middleware stack for Praelatus
var DefaultMiddleware = []func(http.Handler) http.Handler{
	Logger,
}
