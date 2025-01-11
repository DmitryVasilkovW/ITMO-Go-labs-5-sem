package client

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"gitlab.com/slon/shad-go/coverme/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Add(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/todo/create", r.URL.Path)

		w.WriteHeader(http.StatusCreated)
		err := json.NewEncoder(w).Encode(&models.Todo{ID: 1, Title: "Test Todo"})
		if err != nil {
			return
		}
	}))
	defer mockServer.Close()

	client := New(mockServer.URL)
	req := &models.AddRequest{Title: "Test Todo", Content: "Test Content"}

	todo, err := client.Add(req)
	require.NoError(t, err)
	require.NotNil(t, todo)
	require.Equal(t, "Test Todo", todo.Title)
}

func TestClient_Add_ErrorStatus(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer mockServer.Close()

	client := New(mockServer.URL)
	req := &models.AddRequest{Title: "Test Todo", Content: "Test Content"}

	todo, err := client.Add(req)
	require.Error(t, err)
	require.Nil(t, todo)
	require.Contains(t, err.Error(), "unexpected status code")
}

func TestClient_Get(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/todo/1", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(&models.Todo{ID: 1, Title: "Test Todo"})
		if err != nil {
			return
		}
	}))
	defer mockServer.Close()

	client := New(mockServer.URL)
	todo, err := client.Get(1)
	require.NoError(t, err)
	require.NotNil(t, todo)
	require.Equal(t, "Test Todo", todo.Title)
}

func TestClient_Get_ErrorStatus(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	client := New(mockServer.URL)
	todo, err := client.Get(1)
	require.Error(t, err)
	require.Nil(t, todo)
	require.Contains(t, err.Error(), "unexpected status code")
}

func TestClient_List(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/todo", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode([]*models.Todo{
			{ID: 1, Title: "Test Todo 1"},
			{ID: 2, Title: "Test Todo 2"},
		})
		if err != nil {
			return
		}
	}))
	defer mockServer.Close()

	client := New(mockServer.URL)
	todos, err := client.List()
	require.NoError(t, err)
	require.Len(t, todos, 2)
	require.Equal(t, "Test Todo 1", todos[0].Title)
}

func TestClient_List_ErrorStatus(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	client := New(mockServer.URL)
	todos, err := client.List()
	require.Error(t, err)
	require.Nil(t, todos)
	require.Contains(t, err.Error(), "unexpected status code")
}

func TestClient_Finish(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/todo/1/finish", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	client := New(mockServer.URL)
	err := client.Finish(1)
	require.NoError(t, err)
}

func TestClient_Finish_ErrorStatus(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	client := New(mockServer.URL)
	err := client.Finish(1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected status code")
}

func TestClient_Add_ErrorPost(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer mockServer.Close()

	client := New(mockServer.URL)
	req := &models.AddRequest{Title: "Test Todo", Content: "Test Content"}

	todo, err := client.Add(req)
	require.Error(t, err)
	require.Nil(t, todo)
}

func TestClient_Add_ErrorDecode(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte("{not_a_valid_json}"))
		if err != nil {
			return
		}
	}))
	defer mockServer.Close()

	client := New(mockServer.URL)
	req := &models.AddRequest{Title: "Test Todo", Content: "Test Content"}

	todo, err := client.Add(req)
	require.Error(t, err)
	require.Nil(t, todo)
}

func TestClient_Get_ErrorGet(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer mockServer.Close()

	client := New(mockServer.URL)
	todo, err := client.Get(1)
	require.Error(t, err)
	require.Nil(t, todo)
}

func TestClient_Get_ErrorDecode(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("{not_a_valid_json}"))
		if err != nil {
			return
		}
	}))
	defer mockServer.Close()

	client := New(mockServer.URL)
	todo, err := client.Get(1)
	require.Error(t, err)
	require.Nil(t, todo)
}

func TestClient_Add_ErrorMarshal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(server.URL)

	addRequest := &models.AddRequest{
		Title:   "Test Todo",
		Content: string([]byte{0x00}),
	}

	todo, err := client.Add(addRequest)

	require.Error(t, err)
	require.Nil(t, todo)
}
