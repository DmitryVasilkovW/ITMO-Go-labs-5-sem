package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"gitlab.com/slon/shad-go/coverme/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/slon/shad-go/coverme/models"
)

func TestInitRoutes(t *testing.T) {
	mockStorage := &mocks.MockStorage{}
	app := New(mockStorage)
	app.initRoutes()

	require.NotNil(t, app.router)
}

func TestStatus(t *testing.T) {
	app := New(nil)
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	app.status(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	body, _ := io.ReadAll(rr.Body)
	require.Contains(t, string(body), "API is up and working!")
}

func TestList(t *testing.T) {
	mockStorage := &mocks.MockStorage{}
	mockStorage.On("GetAll").Return([]models.Todo{{ID: 1, Title: "Test", Content: "Content"}}, nil)

	app := New(mockStorage)
	req, _ := http.NewRequest("GET", "/todo", nil)
	rr := httptest.NewRecorder()

	app.list(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	mockStorage.AssertCalled(t, "GetAll")
}

func TestAddTodo_Success(t *testing.T) {
	mockStorage := &mocks.MockStorage{}
	mockStorage.On("AddTodo", "Test", "Content").Return(models.Todo{ID: 1, Title: "Test", Content: "Content"}, nil)

	app := New(mockStorage)
	body, _ := json.Marshal(map[string]string{"title": "Test", "content": "Content"})
	req, _ := http.NewRequest("POST", "/todo/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.addTodo(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	mockStorage.AssertCalled(t, "AddTodo", "Test", "Content")
}

func TestAddTodo_EmptyTitle(t *testing.T) {
	mockStorage := &mocks.MockStorage{}
	app := New(mockStorage)

	body, _ := json.Marshal(map[string]string{"title": "", "content": "Content"})
	req, _ := http.NewRequest("POST", "/todo/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.addTodo(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetTodo_Success(t *testing.T) {
	mockStorage := &mocks.MockStorage{}
	mockStorage.On("GetTodo", models.ID(1)).Return(models.Todo{ID: 1, Title: "Test", Content: "Content"}, nil)

	app := New(mockStorage)
	req, _ := http.NewRequest("GET", "/todo/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	app.getTodo(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	mockStorage.AssertCalled(t, "GetTodo", models.ID(1))
}

func TestGetTodo_NotFound(t *testing.T) {
	mockStorage := &mocks.MockStorage{}
	mockStorage.On("GetTodo", models.ID(1)).Return(models.Todo{}, errors.New("not found"))

	app := New(mockStorage)
	req, _ := http.NewRequest("GET", "/todo/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	app.getTodo(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestFinishTodo_Success(t *testing.T) {
	mockStorage := &mocks.MockStorage{}
	mockStorage.On("FinishTodo", models.ID(1)).Return(nil)

	app := New(mockStorage)
	req, _ := http.NewRequest("POST", "/todo/1/finish", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	app.finishTodo(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	mockStorage.AssertCalled(t, "FinishTodo", models.ID(1))
}

func TestFinishTodo_Error(t *testing.T) {
	mockStorage := &mocks.MockStorage{}
	mockStorage.On("FinishTodo", models.ID(1)).Return(errors.New("db error"))

	app := New(mockStorage)
	req, _ := http.NewRequest("POST", "/todo/1/finish", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	app.finishTodo(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestAddTodo_DecodeError(t *testing.T) {
	mockStorage := &mocks.MockStorage{}
	app := New(mockStorage)

	body := bytes.NewReader([]byte("{invalid-json"))
	req, _ := http.NewRequest("POST", "/todo/create", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.addTodo(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "payload is required")
}

func TestAddTodo_DBError(t *testing.T) {
	mockStorage := &mocks.MockStorage{}
	mockStorage.On("AddTodo", "Test", "Content").Return(models.Todo{}, errors.New("db error"))

	app := New(mockStorage)
	body, _ := json.Marshal(map[string]string{"title": "Test", "content": "Content"})
	req, _ := http.NewRequest("POST", "/todo/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.addTodo(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	mockStorage.AssertCalled(t, "AddTodo", "Test", "Content")
}

func TestList_Error(t *testing.T) {
	mockStorage := &mocks.MockStorage{}
	mockStorage.On("GetAll").Return(nil, errors.New("db error"))

	app := New(mockStorage)
	req, _ := http.NewRequest("GET", "/todo", nil)
	rr := httptest.NewRecorder()

	app.list(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	mockStorage.AssertCalled(t, "GetAll")
}

func TestGetTodo_InvalidID(t *testing.T) {
	mockStorage := &mocks.MockStorage{}
	app := New(mockStorage)

	req, _ := http.NewRequest("GET", "/todo/invalidID", nil)
	rr := httptest.NewRecorder()

	app.getTodo(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "ID must be an int")
}

func TestFinishTodo_InvalidID(t *testing.T) {
	mockStorage := &mocks.MockStorage{}
	app := New(mockStorage)

	req, _ := http.NewRequest("POST", "/todo/invalidID/finish", nil)
	rr := httptest.NewRecorder()

	app.finishTodo(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "ID must be an int")
}
