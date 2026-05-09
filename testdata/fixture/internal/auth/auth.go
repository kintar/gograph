package auth

import (
	"errors"

	"github.com/ozgurcd/gograph/testdata/fixture/internal/db"
)

var ErrUnauthorized = errors.New("unauthorized access")

type Service struct {
	repo db.Repository
}

func NewService(repo db.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Login(id string) (string, error) {
	name, err := s.repo.GetUser(id)
	if err != nil {
		return "", err
	}
	if name == "" {
		return "", ErrUnauthorized
	}
	return name, nil
}
