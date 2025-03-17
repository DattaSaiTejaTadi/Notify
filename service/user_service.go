package service

import (
	"Notify/models"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) RegisterUser(name, domainID, password, role string) (*models.User, error) {
	id := uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	userRole := strings.ToLower(role)
	fmt.Println("role: " + userRole)
	if userRole != "admin" && (userRole != "senior" && userRole != "member") {
		return nil, errors.New("invalid role")
	}

	user := models.User{
		ID:          id,
		Name:        name,
		DomainID:    domainID,
		Password:    string(hashedPassword),
		Role:        role,
		ActiveState: false,
	}

	if err := s.store.CreateUser(user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *service) LoginUser(domainID, password string) (*models.User, error) {
	user, err := s.store.GetUserByDomainID(domainID)
	log.Println(domainID, password)
	log.Println(user)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
func (s *service) GetMembers() ([]models.UserResponse, error) {
	users, err := s.store.GetMembers()
	return users, err
}
func (s *service) GetAllUsers() ([]models.UserResponse, error) {
	users, err := s.store.GetAllUsers()
	return users, err
}
func (s *service) DeleteUser(userID string) error {

	return nil
}

func (s *service) UpdateActivityAsTrue(userID string) error {
	return s.store.UpdateActivityAsTrue(userID)
}
func (s *service) UpdateActivityAsFalse(userID string) error {
	return s.store.UpdateActivityAsFalse(userID)
}
