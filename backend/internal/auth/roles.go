package auth

import "fmt"

type UserRole int

// The list of roles that a user can have. The order of the roles is important.
// If we add a new role, it should be added at the end of the list.
// If we remove a role, we should not remove it from the list, but instead mark it as deprecated.
// This is to ensure that the order of the roles is not changed.
// The order of the roles is important because the roles are stored as integers in the database.
const (
	SrRole UserRole = iota
	PeerNavigatorRole
	AdherenceCounselorRole
	AdminRole
)

// String returns the string representation of the role.
func (ur UserRole) String() string {
	// The order of the roles is important because the roles are stored as integers in the database.
	// They must appear in the same order as they are defined above.
	return [...]string{"SrRole", "PeerNavigatorRole", "AdherenceCounselorRole", "AdminRole"}[ur]
}

// ToUserRole converts a string to a UserRole.
// If the string is not a valid role, an error is returned.
func ToUserRole(role string) (UserRole, error) {
	switch role {
	case "SrRole":
		return SrRole, nil
	case "PeerNavigatorRole":
		return PeerNavigatorRole, nil
	case "AdherenceCounselorRole":
		return AdherenceCounselorRole, nil
	case "AdminRole":
		return AdminRole, nil
	default:
		return -1, fmt.Errorf("%s is an invalid role", role) //nolint: goerr113
	}
}
