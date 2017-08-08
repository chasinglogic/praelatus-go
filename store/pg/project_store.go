package pg

import (
	"database/sql"

	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/store"
)

// ProjectStore contains methods for storing and retrieving Projects from a
// Postgres DB
type ProjectStore struct {
	db *sql.DB
}

func intoProject(row rowScanner, p *models.Project) error {
	err := row.Scan(&p.ID, &p.CreatedDate, &p.Name, &p.Key,
		&p.Repo, &p.Homepage, &p.IconURL, &p.Lead.ID,
		&p.Lead.Username, &p.Lead.Email, &p.Lead.FullName,
		&p.Lead.ProfilePic)
	if err != nil {
		return err
	}

	return err
}

// Get gets a project by it's ID in a postgres DB if the given user
// has the appropriate permissions.
func (ps *ProjectStore) Get(u models.User, p *models.Project) error {
	row := ps.db.QueryRow(`
SELECT p.id, p.created_date, p.name, 
       p.key, p.repo, p.homepage, p.icon_url, 
       lead.id, lead.username, lead.email, lead.full_name,  
       lead.profile_picture
FROM projects  AS p
INNER JOIN users AS lead ON lead.id = p.lead_id
FULL JOIN project_permission_schemes AS 
     project_scheme ON p.id = project_scheme.project_id
LEFT JOIN permission_schemes AS scheme ON scheme.id = project_scheme.permission_scheme_id
LEFT JOIN permission_scheme_permissions AS perms ON perms.scheme_id = project_scheme.permission_scheme_id
LEFT JOIN permissions AS perm ON perm.id = perms.perm_id
LEFT JOIN roles AS r ON perms.role_id = r.id
LEFT JOIN user_roles AS roles ON roles.role_id = perms.role_id
LEFT JOIN users AS u ON roles.user_id = u.id
WHERE (p.id = $2 OR p.key = $3)
AND (
    (perm.name = 'VIEW_PROJECT' AND (roles.user_id = $1 OR r.name = 'Anonymous'))
    OR 
    (select is_admin from users where users.id = $1 and users.is_admin = true)
)
;
`,
		u.ID, p.ID, p.Key)

	err := intoProject(row, p)
	return handlePqErr(err)
}

// GetAll returns all projects that the given user has access to
func (ps *ProjectStore) GetAll(u models.User) ([]models.Project, error) {
	var projects []models.Project

	rows, err := ps.db.Query(`
SELECT p.id, p.created_date, p.name, 
       p.key, p.repo, p.homepage, p.icon_url, 
       lead.id, lead.username, lead.email, lead.full_name,  
       lead.profile_picture
FROM projects AS p
INNER JOIN users AS lead ON p.lead_id = lead.id
FULL JOIN project_permission_schemes AS 
     project_scheme ON p.id = project_scheme.project_id
LEFT JOIN permission_schemes AS scheme ON scheme.id = project_scheme.permission_scheme_id
LEFT JOIN permission_scheme_permissions AS perms ON perms.scheme_id = project_scheme.permission_scheme_id
LEFT JOIN permissions AS perm ON perm.id = perms.perm_id
LEFT JOIN roles AS r ON perms.role_id = r.id
LEFT JOIN user_roles AS roles ON roles.role_id = perms.role_id
LEFT JOIN users AS u ON roles.user_id = u.id
WHERE (perm.name = 'VIEW_PROJECT'
AND (roles.user_id = $1 OR r.name = 'Anonymous'))
OR (select is_admin from users where users.id = $2 and users.is_admin = true);
`,
		u.ID, u.ID)

	if err != nil {
		return projects, handlePqErr(err)
	}

	for rows.Next() {
		var p models.Project

		err = intoProject(rows, &p)
		if err != nil {
			return projects, handlePqErr(err)
		}

		projects = append(projects, p)
	}

	return projects, nil
}

// New creates a new Project in the database.
func (ps *ProjectStore) New(project *models.Project) error {
	err := ps.db.QueryRow(`
INSERT INTO projects (name, key, repo, homepage, icon_url, lead_id) 
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
`,
		project.Name, project.Key, project.Repo, project.Homepage,
		project.IconURL, project.Lead.ID).
		Scan(&project.ID)

	return handlePqErr(err)
}

// Create creates a project in the database if the given user is a
// system admin
func (ps *ProjectStore) Create(u models.User, project *models.Project) error {
	if !checkIfAdmin(ps.db, u.ID) {
		return store.ErrPermissionDenied
	}

	return ps.New(project)
}

// Save updates a Project in the database.
func (ps *ProjectStore) Save(u models.User, project models.Project) error {
	if !checkIfAdmin(ps.db, u.ID) {
		return store.ErrPermissionDenied
	}

	_, err := ps.db.Exec(`
UPDATE projects SET 
(name, key, repo, homepage, icon_url, lead_id) 
= ($1, $2, $3, $4, $5, $6)
WHERE projects.id = $7;
`,
		project.Name, project.Key, project.Repo, project.Homepage,
		project.IconURL, project.Lead.ID, project.ID)

	return handlePqErr(err)
}

// Remove updates a Project in the database.
func (ps *ProjectStore) Remove(u models.User, project models.Project) error {
	if !checkIfAdmin(ps.db, u.ID) {
		return store.ErrPermissionDenied
	}

	tx, err := ps.db.Begin()
	if err != nil {
		return handlePqErr(err)
	}

	_, err = tx.Exec(`
DELETE FROM field_tickettype_project WHERE project_id = $1;
`,
		project.ID)

	if err != nil {
		return handlePqErr(tx.Rollback())
	}

	_, err = tx.Exec(`
DELETE FROM project_permission_schemes WHERE project_id = $1;
`,
		project.ID)

	if err != nil {
		return handlePqErr(tx.Rollback())
	}

	_, err = tx.Exec(`
DELETE FROM user_roles WHERE project_id = $1;
`,
		project.ID)

	if err != nil {
		return handlePqErr(tx.Rollback())
	}

	_, err = tx.Exec(`
DELETE FROM field_values
WHERE ticket_id in(SELECT id FROM tickets WHERE project_id = $1);
`,
		project.ID)

	if err != nil {
		return handlePqErr(tx.Rollback())
	}

	_, err = tx.Exec(`
DELETE FROM tickets_labels 
WHERE ticket_id in(SELECT id FROM tickets WHERE project_id = $1);
`,
		project.ID)

	if err != nil {
		return handlePqErr(tx.Rollback())
	}

	_, err = tx.Exec(`DELETE FROM tickets WHERE project_id = $1;`,
		project.ID)

	if err != nil {
		tx.Rollback()
		return handlePqErr(tx.Rollback())
	}

	_, err = tx.Exec(`DELETE FROM projects WHERE id = $1;`,
		project.ID)

	if err != nil {
		return handlePqErr(tx.Rollback())
	}

	return handlePqErr(tx.Commit())
}

// SetPermissionScheme will associate the given permission scheme with the given project
func (ps *ProjectStore) SetPermissionScheme(u models.User, p models.Project, scheme models.PermissionScheme) error {
	if !checkPermission(ps.db, "ADMIN_PROJECT", p.ID, u.ID) {
		return store.ErrPermissionDenied
	}

	_, err := ps.db.Exec(`
INSERT INTO project_permission_schemes 
(permission_scheme_id, project_id)
VALUES ($1, $2)
`,
		scheme.ID, p.ID)

	return handlePqErr(err)
}
