package middleware

import (
	"errors"
	"net/http"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/server"
	"github.com/dgrijalva/jwt-go"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				http.Error(w, "Auth token cookie is required", http.StatusUnauthorized)
			} else {
				http.Error(w, "Error retrieving auth token cookie", http.StatusUnauthorized)
			}
			return
		}

		tokenString := cookie.Value

		token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte("FILATIK_SECRET_KEY_FOR_TOKEN"), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
			ctx := server.SetStringToContext(r.Context(), server.ContextKeyUserLogin, claims.Subject)
			r = r.WithContext(ctx)
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
