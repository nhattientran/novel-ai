package domain

// User represents a user in the system
type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"` // never expose in JSON
	Role         string `json:"role"` // "creator" or "reader"
}

// IsCreator returns true if user has creator role
func (u *User) IsCreator() bool {
	return u.Role == "creator"
}
