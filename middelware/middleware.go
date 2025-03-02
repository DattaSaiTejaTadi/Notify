package middelware

import (
	"Notify/constants"
	"Notify/models"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWT(handler http.HandlerFunc, options *models.Validations) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(constants.JWT_COOKIE)
		if err != nil {
			http.Error(w, "Cookie missing", http.StatusUnauthorized)
			return
		}
		claims, err := ParseJWT(cookie.Value)
		if err != nil {
			http.Error(w, "invalid token or expired", http.StatusUnauthorized)
			return
		}
		if options != nil && claims["role"] != options.Role {
			fmt.Print(options)
			http.Error(w, "Content not accessible", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), constants.UserDetailsFromToken, claims)
		handler.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GenerateJWT(data map[string]interface{}) (string, error) {

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.

	claims := jwt.MapClaims{
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		"exp": time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 days expiration
		"iat": time.Now().Unix(),                          // Issued at current time
	}

	for key, value := range data {
		claims[key] = value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(constants.HmacSampleSecret)

	return tokenString, err
}
func ParseJWT(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return constants.HmacSampleSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("couldn't parse the token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
