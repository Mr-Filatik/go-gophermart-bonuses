package middleware

import (
	"context"
	"errors"
	"net/http"

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
			ctx := r.Context()
			ctx = context.WithValue(ctx, "login", claims.Subject)
			r = r.WithContext(ctx)
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// var SecretKey = []byte("FILATIK_SECRET_KEY_FOR_TOKEN")

// func AuthMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" {
// 			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
// 			return
// 		}

// 		fmt.Println(authHeader)

// 		parts := strings.Split(authHeader, " ")
// 		if len(parts) != 2 || parts[0] != "Bearer" {
// 			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
// 			return
// 		}

// 		tokenString := parts[1]

// 		// Парсим и валидируем токен
// 		token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
// 			// Проверяем алгоритм подписи
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, errors.New("unexpected signing method")
// 			}
// 			return SecretKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
// 			return
// 		}

// 		// Извлекаем claims из токена
// 		if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
// 			// Сохраняем логин пользователя в контексте запроса
// 			ctx := r.Context()
// 			ctx = context.WithValue(ctx, "login", claims.Subject)
// 			r = r.WithContext(ctx)
// 		} else {
// 			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
// 			return
// 		}

// 		// Передаём управление следующему обработчику
// 		next.ServeHTTP(w, r)
// 	})
// }
