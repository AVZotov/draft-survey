package service

import "github.com/AVZotov/draft-survey/internal/types"

type NoopUserService struct{}

func (s *NoopUserService) Get() (*types.User, error)       { return nil, nil }
func (s *NoopUserService) Save(user *types.User) error     { return nil }
func (s *NoopUserService) SaveSignature(data []byte) error { return nil }
func (s *NoopUserService) Delete() error                   { return nil }
