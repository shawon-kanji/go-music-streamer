package entity

type User struct {
	ID        uint
	Name      string
	Email     string
	Password  string
	UserRoles []string
}

// HasRole checks if the user possesses a specific role
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.UserRoles {
		if role == roleName {
			return true
		}
	}
	return false
}
