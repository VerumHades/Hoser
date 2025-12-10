package user

// User represents a system account internally.
type User struct {
	id           string
	username     string
	passwordHash string
	developer    bool
}
