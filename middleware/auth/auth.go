//go:build !solution

package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type User struct {
	Name  string
	Email string
}

type ctxKey string

var ErrInvalidToken = errors.New("invalid token")

func ContextUser(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(ctxKey("user")).(*User)

	return user, ok
}

type TokenChecker interface {
	CheckToken(ctx context.Context, token string) (*User, error)
}

func CheckAuth(checker TokenChecker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := validateAuthorizationHeader(r.Header.Get("Authorization"))
			if err != nil {
				handleAuthError(w, err)
				return
			}

			token := extractToken(r.Header.Get("Authorization"))
			user, err := validateToken(r.Context(), checker, token)
			if err != nil {
				handleAuthError(w, err)
				return
			}

			ctx := addUserToContext(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func validateAuthorizationHeader(header string) error {
	if header == "" {
		return fmt.Errorf("missing authorization header: %w", ErrInvalidToken)
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return fmt.Errorf("invalid authorization header: %w", ErrInvalidToken)
	}

	return nil
}

func extractToken(header string) string {
	return strings.SplitN(header, " ", 2)[1]
}

func validateToken(ctx context.Context, checker TokenChecker, token string) (interface{}, error) {
	user, err := checker.CheckToken(ctx, token)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			return nil, ErrInvalidToken
		}
		return nil, fmt.Errorf("internal server error: %w", err)
	}
	return user, nil
}

func handleAuthError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrInvalidToken) {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
	} else {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func addUserToContext(ctx context.Context, user interface{}) context.Context {
	return context.WithValue(ctx, ctxKey("user"), user)
}
