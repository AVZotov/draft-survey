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

func (u *userService) Save(incoming *types.User) (*Outcome, error) {
	const op = "UserService.Save"

	existing, err := u.repo.Get()
	if err != nil {
		u.logger.Error(op, err)
		return nil, err
	}

	isNew := existing == nil
	merged := merge(existing, incoming)

	var emailWarning *Outcome

	if errs := u.validator.Validate(merged); len(errs) > 0 {
		for _, fe := range errs {
			if fe.Field == "Email" {
				merged.Email = ""
				u.logger.Warn(op, fe.Field, fe.Message)
				emailWarning = &Outcome{
					Toast: &Notification{
						Kind:    KindWarning,
						Header:  "Invalid Email",
						Message: "Email not saved. All other data updated.",
					},
				}
			}
		}
		if errs := u.validator.Validate(merged); len(errs) > 0 {
			return &Outcome{
				Alert: &Notification{
					Kind:    KindError,
					Header:  "Validation Error",
					Message: errs.Error(),
				},
			}, nil
		}
	}

	if err = u.repo.Save(merged); err != nil {
		u.logger.Error(op, err)
		return nil, err
	}

	u.logger.Info(op, "status:", "user saved successfully")

	if emailWarning != nil {
		return emailWarning, nil
	}

	if isNew {
		return &Outcome{
			Toast: &Notification{
				Kind:    KindSuccess,
				Header:  "Profile Created",
				Message: "You can add extra info or continue",
			},
		}, nil
	}

	return nil, nil
}

func (u *userService) SaveSignature(data []byte) (*Outcome, error) {
	const op = "UserService.SaveSignature"

	if len(data) == 0 {
		u.logger.Warn(op, "empty signature data")
		return nil, nil
	}

	if err := u.repo.SaveSignature(data, ""); err != nil {
		u.logger.Error(op, err)
		return nil, err
	}

	u.logger.Info(op, "signature saved successfully")
	return &Outcome{
		Toast: &Notification{
			Kind:    KindSuccess,
			Header:  "Saved",
			Message: "Signature updated",
		},
	}, nil
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

func merge(existing, incoming *types.User) *types.User {
	if existing == nil {
		return incoming
	}
	if incoming.FirstName != "" {
		existing.FirstName = incoming.FirstName
	}
	if incoming.LastName != "" {
		existing.LastName = incoming.LastName
	}
	if incoming.Position != "" {
		existing.Position = incoming.Position
	}
	if incoming.Email != "" {
		existing.Email = incoming.Email
	}
	if incoming.Company != "" {
		existing.Company = incoming.Company
	}
	if incoming.License != "" {
		existing.License = incoming.License
	}
	if incoming.Country != "" {
		existing.Country = incoming.Country
	}
	if incoming.EmployeeID != "" {
		existing.EmployeeID = incoming.EmployeeID
	}
	return existing
}
