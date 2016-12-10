package routes

import (
	"context"
	"fmt"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

// withClaims checks the request for a valid auth token.
// If valid, the Claims object is added to the request's context
func (env *Env) withClaims(fn http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Auth")
		if err != nil {
			fn.ServeHTTP(w, r)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Make sure token's signature wasn't changed
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected siging method")
			}
			return []byte(os.Getenv("TOKEN_SIGNING_KEY")), nil
		})

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "Claims", *claims)
			fn.ServeHTTP(w, r.WithContext(ctx))
		} else {
			fn.ServeHTTP(w, r)
			return
		}
	})
}

func (env *Env) authRequired(fn http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c := r.Context().Value("Claims"); c != nil {
			if _, ok := c.(Claims); ok {
				fn.ServeHTTP(w, r)
			}
		} else {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}
	})
}
