package store

import "Notify/models"

type Store interface {
	CreateUser(user models.User) error
	GetUserByDomainID(domainID string) (*models.User, error)
	GetMembers() ([]models.UserResponse, error)
	GetAllUsers() ([]models.UserResponse, error)
	UpdateActivityAsTrue(userID string) error
	UpdateActivityAsFalse(userID string) error
}
