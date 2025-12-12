package user

// UserPublicView exposes safe read-only public information about a user.
type UserPublicView struct {
	ID       string
	Username string
	IsDev    bool
}

// UserPrivateView exposes sensitive internal information about a user.
type UserPrivateView struct {
	ID           string
	Username     string
	IsDev        bool
	PasswordHash string
}

// ToPublicView converts a User into a UserPublicView.
func (entity *User) ToPublicView() *UserPublicView {
	publicView := &UserPublicView{
		ID:       entity.id,
		Username: entity.username,
		IsDev:    entity.developer,
	}
	return publicView
}

// ToPrivateView converts a User into a UserPrivateView.
func (entity *User) ToPrivateView() *UserPrivateView {
	privateView := &UserPrivateView{
		ID:           entity.id,
		Username:     entity.username,
		IsDev:        entity.developer,
		PasswordHash: entity.passwordHash,
	}
	return privateView
}
