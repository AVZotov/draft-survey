package service

import (
	"github.com/AVZotov/draft-survey/internal/logger"
	"github.com/AVZotov/draft-survey/internal/storage"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/internal/validation"
)

type userService struct {
	repo      storage.UserRepository
	logger    logger.Logger
	validator validation.Validator
}

func NewUserService(repo storage.UserRepository, log logger.Logger, v validation.Validator) UserService {
	return &userService{
		repo:      repo,
		logger:    log,
		validator: v,
	}
}

func (u *userService) Get() (*types.User, error) {
	const op = "UserService.Get"

	user, err := u.repo.Get()
	if err != nil {
		u.logger.Error(op, err)
		return nil, err
	}

	return user, nil
}

func (u *userService) Save(user *types.User) error {
	const op = "UserService.Save"

	if errs := u.validator.Validate(user); len(errs) > 0 {
		u.logger.Warn("validation failed", "caller:", op, "errors", errs.Error())
		return errs
	}

	if err := u.repo.Save(user); err != nil {
		u.logger.Error(op, err)
		return err
	}

	u.logger.Info(op, "user saved successfully")
	return nil
}

func (u *userService) SaveSignature(data []byte) error {
	const op = "UserService.SaveSignature"

	if len(data) == 0 {
		u.logger.Warn(op, "empty signature data")
		return nil
	}

	if err := u.repo.SaveSignature(data, ""); err != nil {
		u.logger.Error(op, err)
		return err
	}

	u.logger.Info(op, "signature saved successfully")
	return nil
}

func (u *userService) Delete() error {
	const op = "UserService.Delete"

	if err := u.repo.Delete(); err != nil {
		u.logger.Error(op, err)
		return err
	}

	u.logger.Info(op, "user deleted")
	return nil
}
