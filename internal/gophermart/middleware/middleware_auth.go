package middleware

import (
	"errors"
	"net/http"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/server"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddlewareFactory(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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

			token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
				subject, sErr := claims.GetSubject()
				if sErr != nil {
					http.Error(w, "Invalid token claims", http.StatusUnauthorized)
					return
				}
				ctx := server.SetStringToContext(r.Context(), server.ContextKeyUserLogin, subject)
				r = r.WithContext(ctx)
			} else {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
