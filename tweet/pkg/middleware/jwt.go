package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Expecting a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "invalid Authorization format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		_, claims, err := ValidateJwt(token)
		if err != nil {
			http.Error(w, "invalid or expired access token", http.StatusUnauthorized)
			return
		}

		w.Header().Set("usename", claims["name"].(string))
		w.Header().Set("id", claims["ID"].(string))

		// If valid, call the next handler
		next.ServeHTTP(w, r)
	})
}

func ValidateJwt(signedToken string) (*jwt.Token, jwt.MapClaims, error) {
	// Load the RSA public key from a file
	publicKey, err := os.ReadFile("../public.pem")
	if err != nil {
		log.Fatal(err)
	}
	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the signed JWT and verify it with the RSA public key
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexepcted signing method: %v", token.Header["alg"])
		}
		return rsaPublicKey, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return token, claims, nil
	}
	return nil, nil, err
}
