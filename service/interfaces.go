package service

import "Notify/models"

type Service interface {
	RegisterUser(name, domainID, password, role string) (*models.User, error)
	LoginUser(domainID, password string) (*models.User, error)
	GetMembers() ([]models.UserResponse, error)
	GetAllUsers() ([]models.UserResponse, error)
	DeleteUser(userID string) error
	UpdateActivityAsTrue(userID string) error
	UpdateActivityAsFalse(userID string) error
}
