package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

var jwtSecret = []byte(viper.GetString("JWT_SECRET"))

func JWTMiddleware(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request has a valid JWT token
		// If it does, call the next handler
		// If it does not, return a 401 Unauthorized response
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := bearerToken[1]
		claims := &jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("User %s authenticated", (*claims)["username"])
		n.ServeHTTP(w, r)
	})
}

func ValidateJWT(token string) (jwt.MapClaims, error) {
	tokenStr := strings.TrimPrefix(token, "Bearer ")

	claims := jwt.MapClaims{}
	t, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !t.Valid {
		return nil, err
	}

	return claims, nil
}
