package services

import (
	"context"

	"github.com/okiww/billing-loan-system/internal/user/repositories"
	"github.com/okiww/billing-loan-system/pkg/logger"
	"github.com/sirupsen/logrus"
)

type userService struct {
	userRepo repositories.UserRepositoryInterface
}

func (u *userService) UpdateUserToDelinquent(ctx context.Context, userID int32) error {
	logger.GetLogger().Info("[UserService][UpdateUserToDelinquent]")
	err := u.userRepo.UpdateUserToDelinquent(ctx, userID)
	if err != nil {
		logger.GetLogger().WithFields(logrus.Fields{
			"user_id": userID,
		}).Errorf("[UserService][UpdateUserToDelinquent] error update user to delinquent %v", err)
		return err
	}

	return nil
}

func (u *userService) IsDelinquent(ctx context.Context, userID int32) (bool, error) {
	logger.GetLogger().Info("[UserService][IsDelinquent]")
	user, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		logger.GetLogger().WithFields(logrus.Fields{
			"user_id": userID,
		}).Errorf("[UserService][IsDelinquent] error get user by id %v", err)
		return false, err
	}

	return user.IsDelinquent, nil
}

func (u *userService) UpdateUserToNotDelinquent(ctx context.Context, userID int32) error {
	logger.GetLogger().Info("[UserService][UpdateUserToNotDelinquent]")
	err := u.userRepo.UpdateUserToNotDelinquent(ctx, userID)
	if err != nil {
		logger.GetLogger().WithFields(logrus.Fields{
			"user_id": userID,
		}).Errorf("[UserService][UpdateUserToNotDelinquent] error update user to delinquent %v", err)
		return err
	}

	return nil
}

func NewUserService(userRepo repositories.UserRepositoryInterface) UserServiceInterface {
	return &userService{userRepo}
}

type UserServiceInterface interface {
	UpdateUserToDelinquent(ctx context.Context, userID int32) error
	IsDelinquent(ctx context.Context, userID int32) (bool, error)
	UpdateUserToNotDelinquent(ctx context.Context, userID int32) error
}
