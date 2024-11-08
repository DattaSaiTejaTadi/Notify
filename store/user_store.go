package store

import (
	"database/sql"
)

type User struct {
	ID          string
	Name        string
	DomainID    string
	Password    string
	Role        string
	ActiveState bool
}

type UserStore struct {
	DB *sql.DB
}

func (s *UserStore) CreateUser(user User) error {
	query := `INSERT INTO users (id, name, domainid, password,role, active_state) VALUES (?, ?, ?, ?, ?,?)`
	_, err := s.DB.Exec(query, user.ID, user.Name, user.DomainID, user.Password, user.Role, user.ActiveState)
	return err
}

func (s *UserStore) GetUserByDomainID(domainID string) (*User, error) {
	query := `SELECT id, name, domainid, password, active_state FROM users WHERE domainid = ?`
	row := s.DB.QueryRow(query, domainID)

	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.DomainID, &user.Password, &user.ActiveState); err != nil {
		return nil, err
	}
	return &user, nil
}
