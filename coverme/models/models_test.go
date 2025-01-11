package models

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInMemoryStorage_AddTodo(t *testing.T) {
	storage := NewInMemoryStorage()

	todo, err := storage.AddTodo("Test Title", "Test Content")
	require.NoError(t, err)
	require.NotNil(t, todo)
	require.Equal(t, "Test Title", todo.Title)
	require.Equal(t, "Test Content", todo.Content)
	require.Equal(t, false, todo.Finished)
}

func TestInMemoryStorage_GetTodo(t *testing.T) {
	storage := NewInMemoryStorage()
	todo, _ := storage.AddTodo("Test Title", "Test Content")

	retrievedTodo, err := storage.GetTodo(todo.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedTodo)
	require.Equal(t, todo.ID, retrievedTodo.ID)
	require.Equal(t, todo.Title, retrievedTodo.Title)
}

func TestInMemoryStorage_GetTodo_NotFound(t *testing.T) {
	storage := NewInMemoryStorage()

	_, err := storage.GetTodo(1)
	require.Error(t, err)
}

func TestInMemoryStorage_GetAll(t *testing.T) {
	storage := NewInMemoryStorage()

	_, err := storage.AddTodo("Title 1", "Content 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = storage.AddTodo("Title 2", "Content 2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	todos, err := storage.GetAll()
	require.NoError(t, err)
	require.Len(t, todos, 2)
}

func TestInMemoryStorage_FinishTodo(t *testing.T) {
	storage := NewInMemoryStorage()
	todo, _ := storage.AddTodo("Test Title", "Test Content")

	err := storage.FinishTodo(todo.ID)
	require.NoError(t, err)
	require.True(t, todo.Finished)
}

func TestInMemoryStorage_FinishTodo_NotFound(t *testing.T) {
	storage := NewInMemoryStorage()

	err := storage.FinishTodo(1)
	require.Error(t, err)
}

func TestTodo_MarkFinished(t *testing.T) {
	todo := &Todo{
		ID:       1,
		Title:    "Test Todo",
		Content:  "Test Content",
		Finished: false,
	}

	require.False(t, todo.Finished)

	todo.MarkFinished()

	require.True(t, todo.Finished)
}

func TestTodo_MarkUnfinished(t *testing.T) {
	todo := &Todo{
		ID:       1,
		Title:    "Test Todo",
		Content:  "Test Content",
		Finished: true,
	}

	require.True(t, todo.Finished)

	todo.MarkUnfinished()

	require.False(t, todo.Finished)
}

func TestAddRequest_Validation(t *testing.T) {
	req := &AddRequest{
		Title:   "Test Title",
		Content: "Test Content",
	}

	require.Equal(t, "Test Title", req.Title)
	require.Equal(t, "Test Content", req.Content)
}
