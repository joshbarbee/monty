package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	Admin    bool   `json:"admin"`
	Username string `json:"username"`
	jwt.StandardClaims
}

const JWT_EXP = time.Minute * 15

func generateJWT(user *User, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		user.Admin,
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(JWT_EXP).Unix(),
		},
	})

	return token.SignedString(secret)
}

func verifyJWT(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("JWTToken")

			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			cookie_kv := strings.Split(cookie.String(), "=")

			if len(cookie_kv) != 2 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token, err := jwt.ParseWithClaims(cookie_kv[1], &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				return secret, nil
			})

			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
				if time.Now().Unix() > claims.ExpiresAt {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		})
	}
}
