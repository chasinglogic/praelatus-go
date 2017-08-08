package pg

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/store"
)

// TicketStore contains methods for storing and retrieving Tickets from
// Postgres DB
type TicketStore struct {
	db *sql.DB
}

func getOpts(db *sql.DB, fid int64, fo *models.FieldOption) error {
	rows, err := db.Query(`
SELECT option FROM field_options 
WHERE field_id = $1`,
		fid)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var opt string

		err = rows.Scan(&opt)
		if err != nil {
			return err
		}

		fo.Options = append(fo.Options, opt)
	}

	return nil
}

func populateFields(db *sql.DB, t *models.Ticket) error {
	rows, err := db.Query(`
SELECT fv.id, f.name, f.data_type, 
       fv.int_value, fv.flt_value, fv.str_value, 
       fv.opt_value, fv.dte_value, f.id
FROM field_values AS fv
JOIN fields AS f ON f.id = fv.field_id
WHERE fv.ticket_id = $1`,
		t.ID)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		// We need to be able to scan in all the values then determine which
		// actually goes into the model.
		fv := models.FieldValue{}
		var i int
		var f float64
		var s, o string
		var d time.Time
		var fID int64

		err = rows.Scan(&fv.ID, &fv.Name, &fv.DataType, &i, &f, &s, &o, &d, &fID)
		if err != nil {
			return err
		}

		// By Odin's Beard I can't think of a better way to wrangle this mess.
		switch fv.DataType {
		case "FLOAT":
			fv.Value = f
			break
		case "INT":
			fv.Value = i
			break
		case "STRING":
			fv.Value = s
			break
		case "DATE":
			fv.Value = d
			break
		case "OPT":
			fo := &models.FieldOption{}
			fo.Selected = o

			// Fill out the options and defaults.
			e := getOpts(db, fID, fo)
			if e != nil {
				return e
			}

			fv.Value = fo
			break
		default:
			fv.Value = nil
		}

		t.Fields = append(t.Fields, fv)
	}

	return nil
}

func populateTransitions(db *sql.DB, t *models.Ticket) error {
	rows, err := db.Query(`
SELECT t.id, t.name, to_s.id, to_s.name
FROM transitions AS t
JOIN statuses AS to_s ON to_s.id = t.to_status
WHERE t.from_status = $1
AND t.workflow_id = $2`,
		t.Status.ID, t.WorkflowID)
	if err != nil {
		return err
	}

	t.Transitions, err = transitionsFromRows(db, rows)
	return err
}

func populateLabels(db *sql.DB, t *models.Ticket) error {
	rows, err := db.Query(`
SELECT l.id, l.name 
FROM tickets_labels AS tl
JOIN labels AS l ON tl.label_id = l.id
WHERE tl.ticket_id = $1`,
		t.ID)
	if err != nil {
		return err
	}

	for rows.Next() {
		l := models.Label{}

		err = rows.Scan(&l.ID, &l.Name)
		if err != nil {
			return err
		}

		t.Labels = append(t.Labels, l)
	}

	return nil
}

func intoTicket(row rowScanner, db *sql.DB, t *models.Ticket) error {
	err := row.Scan(&t.ID, &t.Key, &t.CreatedDate, &t.UpdatedDate, &t.Summary,
		&t.Description, &t.WorkflowID,
		&t.Assignee.ID, &t.Assignee.Username, &t.Assignee.Email,
		&t.Assignee.FullName, &t.Assignee.ProfilePic,
		&t.Reporter.ID, &t.Reporter.Username, &t.Reporter.Email,
		&t.Reporter.FullName, &t.Reporter.ProfilePic,
		&t.Status.ID, &t.Status.Name,
		&t.Type.ID, &t.Type.Name,
		&t.Project.ID, &t.Project.CreatedDate, &t.Project.Name,
		&t.Project.Key, &t.Project.Homepage, &t.Project.IconURL,
		&t.Project.Repo)
	if err != nil {
		return handlePqErr(err)
	}

	err = populateTransitions(db, t)
	if err != nil {
		return handlePqErr(err)
	}

	err = populateFields(db, t)
	if err != nil {
		return handlePqErr(err)
	}

	err = populateLabels(db, t)
	if err != nil {
		return handlePqErr(err)
	}

	return handlePqErr(err)
}

// Get gets a Ticket from a postgres DB by it's ID
func (ts *TicketStore) Get(u models.User, t *models.Ticket) error {
	row := ts.db.QueryRow(`
SELECT t.id, t.key, t.created_date, 
       t.updated_date, t.summary, t.description, t.workflow_id,
       a.id, a.username, a.email, a.full_name, a.profile_picture, 
       r.id, r.username, r.email, r.full_name, r.profile_picture,
       s.id, s.name,
       tt.id, tt.name,
       p.id, p.created_date, p.name, p.key, p.homepage, p.icon_url, p.repo
FROM tickets AS t 
JOIN users AS a ON a.id = t.assignee_id
JOIN users AS r ON r.id = t.reporter_id
JOIN statuses AS s ON s.id = t.status_id
JOIN ticket_types AS tt ON tt.id = t.ticket_type_id
JOIN projects AS p ON p.id = t.project_id
FULL JOIN project_permission_schemes AS 
     project_scheme ON p.id = project_scheme.project_id
LEFT JOIN permission_schemes AS scheme ON scheme.id = project_scheme.permission_scheme_id
LEFT JOIN permission_scheme_permissions AS perms ON perms.scheme_id = project_scheme.permission_scheme_id
LEFT JOIN permissions AS perm ON perm.id = perms.perm_id
LEFT JOIN roles AS rl ON perms.role_id = rl.id
LEFT JOIN user_roles AS roles ON roles.role_id = perms.role_id
LEFT JOIN users AS u ON roles.user_id = u.id
WHERE (t.id = $1 OR t.key = $2)
AND (
    (perm.name = 'VIEW_PROJECT' AND (roles.user_id = $3 OR rl.name = 'Anonymous'))
    OR 
    (select is_admin from users where users.id = $3 and users.is_admin = true)
)`,
		t.ID, t.Key, u.ID)
	err := intoTicket(row, ts.db, t)
	return handlePqErr(err)
}

// GetAll gets all the Tickets from the database
func (ts *TicketStore) GetAll(u models.User) ([]models.Ticket, error) {
	var tickets []models.Ticket

	rows, err := ts.db.Query(`
SELECT t.id, t.key, t.created_date, 
       t.updated_date, t.summary, t.description, t.workflow_id,
       a.id, a.username, a.email, a.full_name, a.profile_picture, 
       r.id, r.username, r.email, r.full_name, r.profile_picture,
       s.id, s.name,
       tt.id, tt.name,
       p.id, p.created_date, p.name, p.key, p.homepage, p.icon_url, p.repo
FROM tickets AS t 
JOIN users AS a ON a.id = t.assignee_id
JOIN users AS r ON r.id = t.reporter_id
JOIN statuses AS s ON s.id = t.status_id
JOIN ticket_types AS tt ON tt.id = t.ticket_type_id
JOIN projects AS p ON p.id = t.project_id
FULL JOIN project_permission_schemes AS 
     project_scheme ON p.id = project_scheme.project_id
LEFT JOIN permission_schemes AS scheme ON scheme.id = project_scheme.permission_scheme_id
LEFT JOIN permission_scheme_permissions AS perms ON perms.scheme_id = project_scheme.permission_scheme_id
LEFT JOIN permissions AS perm ON perm.id = perms.perm_id
LEFT JOIN roles AS rl ON perms.role_id = rl.id
LEFT JOIN user_roles AS roles ON roles.role_id = perms.role_id
LEFT JOIN users AS u ON roles.user_id = u.id
WHERE (perm.name = 'VIEW_PROJECT'
AND (roles.user_id = $1 OR rl.name = 'Anonymous'))
OR (select is_admin from users where users.id = $1 and users.is_admin = true);
`,
		u.ID)
	if err != nil {
		return tickets, handlePqErr(err)
	}

	for rows.Next() {
		var t models.Ticket

		err = intoTicket(rows, ts.db, &t)
		if err != nil {
			log.Println("Error getting tickets")
			return tickets, handlePqErr(err)
		}

		tickets = append(tickets, t)
	}

	return tickets, nil
}

// GetAllByProject gets all the Tickets from the database based on the given
// project
func (ts *TicketStore) GetAllByProject(u models.User, p models.Project) ([]models.Ticket, error) {
	var tickets []models.Ticket

	if !checkPermission(ts.db, "VIEW_PROJECT", p.ID, u.ID) {
		return nil, store.ErrPermissionDenied
	}

	rows, err := ts.db.Query(`
SELECT t.id, t.key, t.created_date, 
       t.updated_date, t.summary, t.description, t.workflow_id,
       a.id, a.username, a.email, a.full_name, a.profile_picture, 
       r.id, r.username, r.email, r.full_name, r.profile_picture,
       s.id, s.name,
       tt.id, tt.name,
       p.id, p.created_date, p.name, p.key, p.homepage, p.icon_url, p.repo
FROM tickets AS t 
JOIN users AS a ON a.id = t.assignee_id
JOIN users AS r ON r.id = t.reporter_id
JOIN projects AS p ON p.id = t.project_id
JOIN statuses AS s ON s.id = t.status_id
JOIN ticket_types AS tt ON tt.id = t.ticket_type_id
WHERE (p.id = $1
OR p.key = $2);
`,
		p.ID, p.Key)
	if err != nil {
		return tickets, handlePqErr(err)
	}

	for rows.Next() {
		var t models.Ticket

		err = intoTicket(rows, ts.db, &t)
		if err != nil {
			return tickets, handlePqErr(err)
		}

		tickets = append(tickets, t)
	}

	return tickets, nil
}

func saveFieldValue(db *sql.DB, fv models.FieldValue) error {
	// Yay interface{} hacks!  Should probably not abuse 0 values
	// here but I'm not sure that it matters
	fvint, _ := fv.Value.(int)
	fvflt, _ := fv.Value.(float64)
	fvstr, _ := fv.Value.(string)
	fvdte, _ := fv.Value.(time.Time)
	fvopt, _ := fv.Value.(models.FieldOption)

	_, err := db.Exec(`
	UPDATE field_values 
    SET (name, data_type, int_value, flt_value, 
         str_value, dte_value, opt_value) = ($1, $2, $3, $4, $5, $6, $7)
    WHERE id = $8`,
		fv.Name, fv.DataType, fvint, fvflt, fvstr, fvdte, fvopt.Selected, fv.ID)

	return handlePqErr(err)
}

// Save will update an existing ticket in the postgres DB
func (ts *TicketStore) Save(u models.User, project models.Project, ticket models.Ticket) error {
	if !checkPermission(ts.db, "EDIT_TICKET", project.ID, u.ID) {
		return store.ErrPermissionDenied
	}

	_, err := ts.db.Exec(`
    UPDATE tickets SET 
    (summary, description, updated_date) = ($1, $2, $3) 
    WHERE id = $4`,
		ticket.Summary, ticket.Description, time.Now(), ticket.ID)

	for _, fv := range ticket.Fields {
		err = saveFieldValue(ts.db, fv)

		if err != nil {
			return handlePqErr(err)
		}
	}

	return handlePqErr(err)
}

// Remove will update an existing ticket in the postgres DB
func (ts *TicketStore) Remove(u models.User, project models.Project, ticket models.Ticket) error {
	if !checkPermission(ts.db, "DELETE_TICKET", project.ID, u.ID) {
		return store.ErrPermissionDenied
	}

	tx, err := ts.db.Begin()
	if err != nil {
		return handlePqErr(tx.Rollback())
	}

	_, err = tx.Exec(`DELETE FROM field_values WHERE ticket_id = $1;`, ticket.ID)
	if err != nil {
		return handlePqErr(tx.Rollback())
	}

	_, err = tx.Exec(`DELETE FROM tickets_labels WHERE ticket_id = $1;`, ticket.ID)
	if err != nil {
		return handlePqErr(tx.Rollback())
	}

	_, err = tx.Exec(`DELETE FROM tickets WHERE id = $1;`, ticket.ID)
	if err != nil {
		return handlePqErr(tx.Rollback())
	}

	return handlePqErr(tx.Commit())
}

// Create will create the ticket on the given project if the given
// user has the appropriate permission
func (ts *TicketStore) Create(u models.User, project models.Project, ticket *models.Ticket) error {
	if !checkPermission(ts.db, "CREATE_TICKET", project.ID, u.ID) {
		return store.ErrPermissionDenied
	}

	return ts.New(project, ticket)
}

// New will add a new Ticket to the postgres DB
func (ts *TicketStore) New(project models.Project, ticket *models.Ticket) error {
	if project.Key == "" {
		return errors.New("project key required for ticket creation")
	}

	ticket.Key = ts.NextTicketKey(project)

	err := ts.db.QueryRow(`
INSERT INTO tickets 
    (summary, description, project_id, assignee_id, 
     reporter_id, ticket_type_id, status_id, key, 
     workflow_id, created_date, updated_date) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id;
`,
		ticket.Summary, ticket.Description, project.ID,
		ticket.Assignee.ID, ticket.Reporter.ID, ticket.Type.ID,
		ticket.Status.ID, ticket.Key, ticket.WorkflowID,
		ticket.CreatedDate, ticket.UpdatedDate).
		Scan(&ticket.ID)

	if err != nil {
		return handlePqErr(err)
	}

	for _, fv := range ticket.Fields {
		var fieldID int64

		err = ts.db.QueryRow(`SELECT id FROM fields WHERE name = $1`, fv.Name).
			Scan(&fieldID)
		if err != nil {
			return handlePqErr(err)
		}

		// Yay interface{} hacks!
		fvint, _ := fv.Value.(int)
		fvflt, _ := fv.Value.(float64)
		fvstr, _ := fv.Value.(string)
		fvdte, _ := fv.Value.(time.Time)
		fvopt, _ := fv.Value.(models.FieldOption)

		err = ts.db.QueryRow(`
INSERT INTO field_values
    (name, data_type, int_value, flt_value, str_value,
     dte_value, opt_value, ticket_id, field_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id
`,
			fv.Name, fv.DataType, fvint, fvflt, fvstr, fvdte, fvopt.Selected,
			ticket.ID, fieldID).Scan(&fv.ID)
		if err != nil {
			return handlePqErr(err)
		}
	}

	for _, lb := range ticket.Labels {
		_, err = ts.db.Exec(`
INSERT INTO tickets_labels (label_id, ticket_id) 
VALUES ($1, $2)
`,
			lb.ID, ticket.ID)
		if err != nil {
			return handlePqErr(err)
		}
	}

	return handlePqErr(err)
}

// GetComment can be used to get a single comment from the database
func (ts *TicketStore) GetComment(u models.User, cm *models.Comment) error {
	var projectID int64
	var c models.Comment

	err := ts.db.QueryRow(`
SELECT c.id, c.created_date, c.updated_date, c.body, t.key,
       a.id, a.username, a.email, a.full_name, 
       a.profile_picture, p.id
FROM comments AS c
JOIN tickets AS t ON t.id = c.ticket_id
JOIN projects AS p ON p.id = t.project_id
WHERE c.id = $1
`,
		cm.ID).
		Scan(&c.ID, &c.CreatedDate, &c.UpdatedDate,
			&c.Body, &c.TicketKey, &c.Author.ID,
			&c.Author.Username, &c.Author.Email,
			&c.Author.FullName, &c.Author.ProfilePic, projectID)
	if err != nil {
		return err
	}

	if checkPermission(ts.db, "VIEW_PROJECT", projectID, u.ID) {
		return nil
	}

	return store.ErrPermissionDenied
}

// GetComments will return all comments for a ticket based on it's ID
func (ts *TicketStore) GetComments(u models.User, p models.Project, t models.Ticket) ([]models.Comment, error) {
	if !checkPermission(ts.db, "VIEW_PROJECT", u.ID, p.ID) {
		return nil, store.ErrPermissionDenied
	}

	var comments []models.Comment

	rows, err := ts.db.Query(`
SELECT c.id, c.created_date, c.updated_date, c.body, t.key,
       a.id, a.username, a.email, a.full_name, 
       a.profile_picture
FROM comments AS c
JOIN tickets AS t ON t.id = c.ticket_id
JOIN users AS a ON a.id = c.author_id
WHERE t.id = $1
OR t.key = $2
`,
		t.ID, t.Key)
	if err != nil {
		return comments, handlePqErr(err)
	}

	for rows.Next() {
		var c models.Comment

		err := rows.Scan(&c.ID, &c.CreatedDate, &c.UpdatedDate,
			&c.Body, &c.TicketKey, &c.Author.ID, &c.Author.Username,
			&c.Author.Email, &c.Author.FullName, &c.Author.ProfilePic)
		if err != nil {
			return comments, handlePqErr(err)
		}

		comments = append(comments, c)
	}

	return comments, nil
}

// CreateComment will add a new comment to the ticket if the given
// user has the appropriate permissions
func (ts *TicketStore) CreateComment(u models.User, p models.Project, t models.Ticket, c *models.Comment) error {
	if !checkPermission(ts.db, "COMMENT_TICKET", p.ID, u.ID) {
		return store.ErrPermissionDenied
	}

	return ts.NewComment(t, c)
}

// NewComment will add a new Comment to the postgres DB
func (ts *TicketStore) NewComment(t models.Ticket, c *models.Comment) error {
	_, err := ts.db.Exec(`
UPDATE tickets SET updated_date = $1 
WHERE id = $2;`,
		time.Now(), t.ID)
	if err != nil {
		return handlePqErr(err)
	}

	err = ts.db.QueryRow(`
INSERT INTO comments 
    (body, ticket_id, author_id) 
VALUES ($1, $2, $3)
RETURNING id;
`,
		c.Body, t.ID, c.Author.ID).
		Scan(&c.ID)

	return handlePqErr(err)
}

// SaveComment will add a new Comment to the postgres DB
func (ts *TicketStore) SaveComment(u models.User, p models.Project, c models.Comment) error {
	hasPerm := checkPermission(ts.db, "EDIT_COMMENT", p.ID, u.ID)

	if !hasPerm {
		hasPerm = checkPermission(ts.db, "EDIT_OWN_COMMENT", p.ID, u.ID)
		if hasPerm {
			var id *int64

			err := ts.db.
				QueryRow(`
SELECT id FROM comments WHERE id = $1 AND author_id = $2
`,
					c.ID, u.ID).
				Scan(&id)

			if err != nil || id == nil {
				return store.ErrPermissionDenied
			}
		}
	}

	if hasPerm {
		_, err := ts.db.Exec(`
UPDATE comments 
SET (body, updated_date, author_id) = ($1, $2, $3)
WHERE id = $4
`,
			c.Body, time.Now(), c.Author.ID, c.ID)
		return handlePqErr(err)
	}

	return store.ErrPermissionDenied
}

// RemoveComment will add a new Comment to the postgres DB
func (ts *TicketStore) RemoveComment(u models.User, p models.Project, c models.Comment) error {
	hasPerm := checkPermission(ts.db, "REMOVE_COMMENT", p.ID, u.ID)

	if !hasPerm {
		hasPerm = checkPermission(ts.db, "REMOVE_OWN_COMMENT", p.ID, u.ID)
		if hasPerm {
			var id *int64

			err := ts.db.
				QueryRow(`
SELECT id FROM comments WHERE id = $1 AND author_id = $2
`,
					c.ID, u.ID).
				Scan(&id)

			if err != nil || id == nil {
				return store.ErrPermissionDenied
			}
		}
	}

	if hasPerm {
		_, err := ts.db.Exec("DELETE FROM comments WHERE id = $1", c.ID)
		return handlePqErr(err)
	}

	return store.ErrPermissionDenied
}

// NextTicketKey will generate the appropriate number for a ticket key
func (ts *TicketStore) NextTicketKey(p models.Project) string {
	var count int

	err := ts.db.QueryRow(`
    SELECT COUNT(t.id) 
    FROM tickets AS t
    JOIN projects AS p ON p.id = t.project_id
    WHERE p.key = $1`, p.Key).Scan(&count)
	if err != nil {
		handlePqErr(err)
		return ""
	}

	return p.Key + "-" + strconv.Itoa(count+1)
}

// ExecuteTransition will take the given transition and perform it on the given
// ticket
func (ts *TicketStore) ExecuteTransition(u models.User, project models.Project, t *models.Ticket, tr models.Transition) error {
	if !checkPermission(ts.db, "TRANSITION_TICKET", project.ID, u.ID) {
		return store.ErrPermissionDenied
	}

	t.Status = tr.ToStatus

	_, err := ts.db.Exec(`
UPDATE tickets
SET status_id = $1;`, t.Status.ID)

	if err != nil {
		return err
	}

	err = populateTransitions(ts.db, t)
	return handlePqErr(err)
}
