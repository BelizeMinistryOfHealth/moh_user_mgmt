package handlers

// UserResponse is the response object for a user
type UserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Enabled   bool   `json:"enabled"`
	Org       string `json:"org"`
	Role      string `json:"role"`
}
