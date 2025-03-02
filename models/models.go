package models

type User struct {
	ID          string
	Name        string
	DomainID    string
	Password    string
	Role        string
	ActiveState bool
}
type UserResponse struct {
	ID          string
	Name        string
	DomainID    string
	Role        string
	ActiveState bool
}

type Validations struct{
	Role string
}
