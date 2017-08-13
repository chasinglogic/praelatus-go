// Package middleware contains the HTTP middleware used in the api as
// well as utility functions for interacting with them
package middleware

import (
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/praelatus/backend/config"
	"github.com/praelatus/backend/models"
)

type SessionCache interface {
	Get(id string) (models.Session, error)
	Set(id string, user models.Session) error
	Remove(id string) error
}

type MongoCache struct {
	conn *mgo.Session
}

func (m MongoCache) Get(id string) (models.Session, error) {
	var s models.Session

	err := m.conn.DB(config.DBName()).C("sessions").FindId(id).One(&s)
	if err != nil {
		if err.Error() != "not found" {
			log.Println("Unexpected Error:", err)
		}

		return models.Session{}, err
	}

	return s, nil
}

func (m MongoCache) Set(id string, sess models.Session) error {
	sess.ID = id

	err := m.conn.DB(config.DBName()).C("sessions").Insert(&sess)
	return err
}

func (m MongoCache) Remove(id string) error {
	return m.conn.DB(config.DBName()).C("sessios").RemoveId(id)
}

// NewMongoCache returns a session store using MongoDB as the backend.
func NewMongoCache(c *mgo.Session) MongoCache {
	return MongoCache{c}
}

// Cache is the global SessionCache
var Cache SessionCache

func headers(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}

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
	headers,
	Logger,
}
