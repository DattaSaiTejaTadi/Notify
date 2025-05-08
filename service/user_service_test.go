package service_test

import (
	"Notify/models"
	"Notify/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockStore struct {
	CreateUserStub error
}

func (ms *MockStore) CreateUser(user models.User) error {
	return ms.CreateUserStub
}
func (ms *MockStore) GetUserByDomainID(domainID string) (*models.User, error) {
	return nil, nil
}
func (ms *MockStore) GetMembers() ([]models.UserResponse, error) {
	return nil, nil
}
func (ms *MockStore) GetAllUsers() ([]models.UserResponse, error) {
	return nil, nil
}
func (ms *MockStore) UpdateActivityAsTrue(userID string) error { return nil }
func (ms *MockStore) UpdateActivityAsFalse(userID string) error {
	return nil
}

type Input struct {
	Name     string
	DomainId string
	Password string
	Role     string
}
type RegisterTestcase struct {
	Title        string
	InputData    Input
	ExpectedData models.User
	ExpectedErr  error
	MockedErr    error
}

func TestRegisterUser(t *testing.T) {
	var testcases = []RegisterTestcase{
		{
			Title:        "Verify Demo check",
			InputData:    Input{Name: "Teja", DomainId: "teja@gamil.com", Password: "Teja@1234", Role: "admin"},
			ExpectedData: models.User{},
			ExpectedErr:  nil,
			MockedErr:    nil,
		},
	}
	for _, tc := range testcases {
		mocksevice := service.New(&MockStore{CreateUserStub: tc.MockedErr})
		
		actualData, actualerr := mocksevice.RegisterUser(tc.InputData.Name, tc.InputData.DomainId, tc.InputData.Password, tc.InputData.Role)

		assert.Equal(t, tc.ExpectedData, actualData, "checking Data")
		assert.Equal(t, tc.ExpectedErr, actualerr, "checking err")
	}

}
