// Package permission is used to store the various permission in
// praelatus, it acts as a pseudo-enumeration
package permission

// Permission is an alias type to make it's use more clear inside of models.
type Permission string

// These are the permissions available in Praelatus
const (
	ViewProject      Permission = "VIEW_PROJECT"
	AdminProject                = "ADMIN_PROJECT"
	CreateTicket                = "CREATE_TICKET"
	CommentTicket               = "COMMENT_TICKET"
	RemoveComment               = "REMOVE_COMMENT"
	RemoveOwnComment            = "REMOVE_OWN_COMMENT"
	EditOwnComment              = "EDIT_OWN_COMMENT"
	EditComment                 = "EDIT_COMMENT"
	TransitionTicket            = "TRANSITION_TICKET"
	EditTicket                  = "EDIT_TICKET"
	RemoveTicket                = "REMOVE_TICKET"
)

// Permissions holds available permissions in a slice. This is valuable for
// various areas where we need to return all permissions or iterate
// permissions.
var Permissions = [...]Permission{
	ViewProject,
	AdminProject,
	CreateTicket,
	CommentTicket,
	RemoveComment,
	RemoveOwnComment,
	EditOwnComment,
	EditComment,
	TransitionTicket,
	EditTicket,
	RemoveTicket,
}

// ValidPermission will verify that a given permission string is valid.
func ValidPermission(permName Permission) bool {
	for _, p := range Permissions {
		if p == permName {
			return true
		}
	}

	return false
}
