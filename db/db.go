// Package db abstracts database interaction away from the handlers to make for
// easier testing.
package db

import mgo "gopkg.in/mgo.v2"
import "github.com/praelatus/backend/config"

// Connection is an interface used for abstracting away mgo.Sessions
type Connection interface {
	DB(dbName string) DB
}

// DB is an interface used for abstracting away mgo.Database
type DB interface {
	C(collectionName string) Collection
}

// Collection is used to abstract away mgo.Collection
type Collection interface {
	Pipe(pipline interface{}) *mgo.Pipe

	Find(query interface{}) *mgo.Query
	FindId(id interface{}) *mgo.Query

	Count() (n int, err error)

	Insert(docs ...interface{}) error

	Remove(selector interface{}) error
	RemoveId(id interface{}) error
	RemoveAll(selector interface{}) (*mgo.ChangeInfo, error)

	Update(selector interface{}, update interface{}) error
	UpdateId(id interface{}, update interface{}) error
}

// Connect connects to the MongoDB specified by the config file.
func Connect() *mgo.Session {
	session, err := mgo.Dial(config.DBURL())
	if err != nil {
		panic(err)
	}

	return session
}
