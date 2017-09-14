package mongo

import (
	"log"

	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/repo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Cache is a repo.Cache backed by MongoDB
type Cache struct {
	conn *mgo.Session
}

func (m Cache) sessions() *mgo.Collection {
	return m.conn.DB(dbName).C("sessions")
}

func (m Cache) cache() *mgo.Collection {
	return m.conn.DB(dbName).C("cache")
}

// TODO: General cache functions

func (m Cache) Get(key string) (interface{}, error) {
	data := bson.M{}
	err := m.cache().FindId(key).One(&data)
	return data["data"], mongoErr(err)
}

func (m Cache) Set(key string, value interface{}) error {
	data := bson.M{
		"_id":  key,
		"data": value,
	}

	return m.cache().Insert(data)
}

func (m Cache) Remove(key string) error {
	return m.cache().RemoveId(key)
}

// GetSession will return the session with the given id (Token)
func (m Cache) GetSession(id string) (models.Session, error) {
	var s models.Session

	err := m.sessions().FindId(id).One(&s)
	if err != nil {
		if err.Error() != "not found" {
			log.Println("[MONGO_CACHE] ERROR:", err)
		}

		return models.Session{}, err
	}

	return s, nil
}

// SetSession will store a session with the given id (Token)
func (m Cache) SetSession(id string, sess models.Session) error {
	sess.ID = id

	err := m.sessions().Insert(&sess)
	return err
}

// RemoveSession will remove the session with the given id (Token)
func (m Cache) RemoveSession(id string) error {
	return mongoErr(m.sessions().RemoveId(id))
}

// NewCache returns a repo.Cache using MongoDB as the backend.
func NewCache(connURL string) repo.Cache {
	conn, err := mgo.Dial(connURL)
	if err != nil {
		panic(err)
	}

	return Cache{conn}
}
