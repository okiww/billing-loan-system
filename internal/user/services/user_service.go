package services

import "context"

type userService struct{}

func (u *userService) UpdateUserToDelinquent(ctx context.Context, userID int32) error {
	//TODO implement me
	panic("implement me")
}

func (u *userService) IsDelinquent(ctx context.Context, userID int32) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserService() UserServiceInterface {
	return &userService{}
}

type UserServiceInterface interface {
	UpdateUserToDelinquent(ctx context.Context, userID int32) error
	IsDelinquent(ctx context.Context, userID int32) (bool, error)
}
