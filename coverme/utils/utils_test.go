package utils

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type NonSerializable struct {
	Channel chan int
}

func TestRespondJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	data := map[string]string{"message": "success"}

	err := RespondJSON(rec, http.StatusOK, data)

	require.NoError(t, err)

	require.Equal(t, http.StatusOK, rec.Code)

	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	expectedResponse := `{"message":"success"}`
	require.JSONEq(t, expectedResponse, rec.Body.String())
}

func TestServerError(t *testing.T) {
	rec := httptest.NewRecorder()

	ServerError(rec)

	require.Equal(t, http.StatusInternalServerError, rec.Code)

	expectedResponse := "Server encountered an error."
	require.Equal(t, expectedResponse, rec.Body.String())
}

func TestBadRequest(t *testing.T) {
	rec := httptest.NewRecorder()

	message := "Invalid request"
	BadRequest(rec, message)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	require.Equal(t, message, rec.Body.String())
}

func TestRespondJSON_ErrorMarshal(t *testing.T) {
	rec := httptest.NewRecorder()

	data := NonSerializable{
		Channel: make(chan int),
	}

	err := RespondJSON(rec, http.StatusOK, data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "json: unsupported type: chan int")
}
