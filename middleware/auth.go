package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware protects routes that require authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("AuthMiddleware: Checking for token...")
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Println("AuthMiddleware: No token cookie found. Redirecting to login.")
				http.Redirect(w, r, "/login?return_to="+r.URL.Path, http.StatusFound)
				return
			}
			log.Printf("AuthMiddleware: Error getting cookie: %v\n", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		tokenStr := cookie.Value
		log.Printf("AuthMiddleware: Token string: %s\n", tokenStr)
		claims := &jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil // Use secret from config
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				log.Println("AuthMiddleware: Invalid token signature.")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			log.Printf("AuthMiddleware: Error parsing token: %v\n", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if !token.Valid {
			log.Println("AuthMiddleware: Token is not valid.")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user ID to context
		sub := (*claims)["sub"].(float64)
		ctx := context.WithValue(r.Context(), "userID", uint(sub))
		log.Printf("AuthMiddleware: User authenticated. UserID: %v\n", uint(sub))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
