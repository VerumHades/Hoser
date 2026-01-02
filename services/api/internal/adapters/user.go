package adapters

import (
	"common/pkg/domain/user"
	"common/pkg/shared"
	"context"
)

// developerCheckAdapter adapts a UserQueryRepository to the userService interface.
type DeveloperCheckAdapter struct {
	userRepository user.UserQueryRepository
}

// NewDeveloperCheckAdapter constructs a new adapter.
func NewDeveloperCheckAdapter(repo user.UserQueryRepository) *DeveloperCheckAdapter {
	return &DeveloperCheckAdapter{
		userRepository: repo,
	}
}

func (a *DeveloperCheckAdapter) IsUserDeveloper(ctx context.Context, userID shared.UserID) (bool, error) {
	userEntity, err := a.userRepository.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}

	return userEntity.IsDeveloper(), nil
}
