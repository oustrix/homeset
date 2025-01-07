package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"slices"

	"github.com/golang-jwt/jwt/v5"
)

var noAuthURL = []string{
	"/login",
}

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if slices.Contains(noAuthURL, r.URL.Path) {
				next.ServeHTTP(w, r)

				return
			}

			jwtTokenRaw, err := r.Cookie("Authorization")
			if err != nil {
				redirect(w, r, "/login")

				return
			}

			token, err := jwt.Parse(jwtTokenRaw.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(secret), nil
			})
			if err != nil {
				slog.Error(err.Error())
				redirect(w, r, "/login")

				return
			}

			if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
				redirect(w, r, "/login")

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
