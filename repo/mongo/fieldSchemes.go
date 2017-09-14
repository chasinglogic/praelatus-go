package mongo

import (
	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/repo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type fieldSchemeRepo struct {
	conn *mgo.Session
}

func (fs fieldSchemeRepo) coll() *mgo.Collection {
	return fs.conn.DB(dbName).C(fieldSchemes)
}

func (fs fieldSchemeRepo) Get(u *models.User, uid string) (models.FieldScheme, error) {
	if u == nil {
		return models.FieldScheme{}, repo.ErrLoginRequired
	}

	var f models.FieldScheme
	err := fs.coll().FindId(bson.ObjectIdHex(uid)).One(&f)

	return f, mongoErr(err)
}

func (fs fieldSchemeRepo) Update(u *models.User, uid string, updated models.FieldScheme) error {
	if u == nil || !u.IsAdmin {
		return repo.ErrAdminRequired
	}

	// FIXME: Handle what to do with tickets and projects associated with this
	// field scheme

	return mongoErr(fs.coll().UpdateId(bson.ObjectIdHex(uid), updated))
}

func (fs fieldSchemeRepo) Create(u *models.User, fieldScheme models.FieldScheme) (models.FieldScheme, error) {
	if u == nil || !u.IsAdmin {
		return models.FieldScheme{}, repo.ErrAdminRequired
	}

	fieldScheme.ID = bson.NewObjectId()

	err := fs.coll().Insert(fieldScheme)
	return fieldScheme, mongoErr(err)
}

func (fs fieldSchemeRepo) Delete(u *models.User, uid string) error {
	if u == nil || !u.IsAdmin {
		return repo.ErrAdminRequired
	}

	// FIXME: Handle what to do with tickets and projects associated with this
	// workflow

	return mongoErr(fs.coll().RemoveId(bson.ObjectIdHex(uid)))
}

func (fs fieldSchemeRepo) Search(u *models.User, query string) ([]models.FieldScheme, error) {
	if u == nil || !u.IsAdmin {
		return []models.FieldScheme{}, repo.ErrAdminRequired
	}

	var schemes []models.FieldScheme

	q := bson.M{}
	if query != "" {
		q = bson.M{"name": bson.M{"$regex": query, "$options": "i"}}
	}

	err := fs.coll().Find(q).All(&schemes)
	return schemes, mongoErr(err)
}
