package pg

import (
	"database/sql"
	"errors"

	"github.com/praelatus/praelatus/models"
)

// TypeStore is used to store ticket types in a postgres database
type TypeStore struct {
	db *sql.DB
}

// Get will get a ticket type by either name or id whichver is provided in tt
func (ts *TypeStore) Get(tt *models.TicketType) error {
	row := ts.db.QueryRow(`SELECT tt.id, tt.name 
								FROM ticket_types AS tt
								WHERE tt.id = $1
								OR tt.name = $2`, tt.ID, tt.Name)
	return handlePqErr(row.Scan(&tt.ID, &tt.Name))
}

// GetAll will return all ticket types from the database
func (ts *TypeStore) GetAll() ([]models.TicketType, error) {
	var typs []models.TicketType

	rows, err := ts.db.Query(`SELECT tt.id, tt.name 
							  FROM ticket_types AS tt`)
	if err != nil {
		return typs, handlePqErr(err)
	}

	for rows.Next() {
		var tt models.TicketType

		err = rows.Scan(&tt.ID, &tt.Name)
		if err != nil {
			return typs, handlePqErr(err)
		}

		typs = append(typs, tt)
	}

	return typs, nil
}

// New will add a new TicketType to the postgres DB
func (ts *TypeStore) New(tt *models.TicketType) error {
	row := ts.db.QueryRow(`INSERT INTO ticket_types (name) 
						   VALUES ($1)
						   RETURNING id;`, tt.Name)
	return handlePqErr(row.Scan(&tt.ID))
}

// Save will add a new TicketType to the postgres DB
func (ts *TypeStore) Save(tt models.TicketType) error {
	_, err := ts.db.Exec(`UPDATE ticket_types 
						  SET (name) = ($1)
						  WHERE id = $2`, tt.Name, tt.ID)
	return handlePqErr(err)
}

// Remove remoevs a ticket type from the database.
func (ts *TypeStore) Remove(tt models.TicketType) error {
	var c int

	tx, err := ts.db.Begin()
	if err != nil {
		return handlePqErr(err)
	}

	err = tx.QueryRow(`SELECT COUNT(id) FROM tickets
					   WHERE ticket_type_id = $1`, tt.ID).Scan(&c)
	if err != nil {
		tx.Rollback()
		return handlePqErr(err)
	}

	if c > 0 {
		tx.Rollback()
		return errors.New("that type is currently in use, refusing to delete")
	}

	_, err = tx.Exec("DELETE FROM ticket_types WHERE id = $1", tt.ID)
	if err != nil {
		tx.Rollback()
		return handlePqErr(err)
	}

	return handlePqErr(tx.Commit())
}
