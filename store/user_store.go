package store

import (
	"Notify/models"
	"log"
)

func (s *store) CreateUser(user models.User) error {
	query := `INSERT INTO users (id, name, domainid, password,role, active_state) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.DB.Exec(query, user.ID, user.Name, user.DomainID, user.Password, user.Role, user.ActiveState)
	return err
}

func (s *store) ResetUserActiveStates() error {
	query := `UPDATE users set active_state=false`
	_, err := s.DB.Exec(query)
	return err
}

func (s *store) GetUserByDomainID(domainID string) (*models.User, error) {
	query := `SELECT id, name, domainid, password,role, active_state FROM users WHERE domainid = $1`
	row := s.DB.QueryRow(query, domainID)
	//log.Println(row.Err().Error())
	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.DomainID, &user.Password, &user.Role, &user.ActiveState); err != nil {
		log.Println("issue in scanning")
		log.Println(err.Error())
		return nil, err
	}
	return &user, nil
}
func (s *store) GetMembers() ([]models.UserResponse, error) {
	query := `SELECT id, name, domainid,role, active_state FROM users Where role != $1`
	rows, err := s.DB.Query(query, "admin")
	if err != nil {
		return nil, err
	}
	var users []models.UserResponse
	for rows.Next() {
		var user models.UserResponse
		if err := rows.Scan(&user.ID, &user.Name, &user.DomainID, &user.Role, &user.ActiveState); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *store) GetAllUsers() ([]models.UserResponse, error) {
	query := `SELECT id, name, domainid,role, active_state FROM users Where role != $1`
	rows, err := s.DB.Query(query, "admin")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserResponse

	for rows.Next() {
		var user models.UserResponse
		if err := rows.Scan(&user.ID, &user.Name, &user.DomainID, &user.Role, &user.ActiveState); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *store) UpdateActivityAsTrue(userID string) error {
	query := `UPDATE users set active_state=true where id = $1`
	_, err := s.DB.Query(query, userID)
	if err != nil {
		return err
	}
	return nil
}
func (s *store) UpdateActivityAsFalse(userID string) error {
	query := `UPDATE users set active_state=false where id = $1`
	_, err := s.DB.Query(query, userID)
	if err != nil {
		return err
	}
	return nil
}
