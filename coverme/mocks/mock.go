package mocks

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/slon/shad-go/coverme/models"
)

type MockStorage struct {
	mock.Mock
}

// AddTodo - мок для метода AddTodo
func (m *MockStorage) AddTodo(title, content string) (*models.Todo, error) {
	args := m.Called(title, content)
	todo, _ := args.Get(0).(*models.Todo)
	return todo, args.Error(1)
}

// GetTodo - мок для метода GetTodo
func (m *MockStorage) GetTodo(id models.ID) (*models.Todo, error) {
	args := m.Called(id)
	todo, _ := args.Get(0).(*models.Todo)
	return todo, args.Error(1)
}

// GetAll - мок для метода GetAll
func (m *MockStorage) GetAll() ([]*models.Todo, error) {
	args := m.Called()
	todos, _ := args.Get(0).([]*models.Todo)
	return todos, args.Error(1)
}

// FinishTodo - мок для метода FinishTodo
func (m *MockStorage) FinishTodo(id models.ID) error {
	args := m.Called(id)
	return args.Error(0)
}
