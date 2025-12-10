package user

// UserView exposes a read-only representation of a user.
type UserView struct {
	ID       string
	Username string
	IsDev    bool
}
