package service

import (
	"Notify/store"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Store *store.UserStore
}

func (s *UserService) RegisterUser(name, domainID, password, role string) (*store.User, error) {
	id := uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := store.User{
		ID:          id,
		Name:        name,
		DomainID:    domainID,
		Password:    string(hashedPassword),
		Role:        role,
		ActiveState: false,
	}

	if err := s.Store.CreateUser(user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) LoginUser(domainID, password string) (*store.User, error) {
	user, err := s.Store.GetUserByDomainID(domainID)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
