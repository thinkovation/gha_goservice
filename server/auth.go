package server

import (
	"context"
	"fmt"
	"goservice/users"

	"net/http"
)

type err struct {
	message string
	code    int
}

func WithError(w http.ResponseWriter, e err) {
	errorResponse := Response{
		Status:  "error",
		Message: e.message,
	}
	writeJSON(w, e.code, errorResponse)
}
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getTokenFromRequest(r)
		if token == "" {
			WithError(w, err{message: "no token", code: 401})
			return
		}

		user := users.GetUserFromToken(token)
		fmt.Println(user)

		// Store the user record in the context
		ctx := context.WithValue(r.Context(), users.UserContextKey, user)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getTokenFromRequest fetches the token sent with a request
func getTokenFromRequest(r *http.Request) string {
	var token string
	token = r.Header.Get("Bearer")
	if token != "" {
		return token
	}
	token = r.URL.Query().Get("token")
	if token != "" {
		return token
	}
	return ""
}
