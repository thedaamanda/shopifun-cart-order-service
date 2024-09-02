package middleware

import (
	"cart-order-service/util/helper/jwt"
	"context"
	"encoding/json"
	"net/http"
)

const emailKey = "email"

func SetUserID(ctx context.Context, email string) context.Context {
	ctx = context.WithValue(ctx, emailKey, email)
	return ctx
}

func GetUserID(ctx context.Context) string {
	userID, ok := ctx.Value(emailKey).(string)
	if !ok {
		return ""
	}
	return userID
}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"Message": "Unauthorized",
				"Data":    nil,
			})
			return
		}

		tokenString = tokenString[len("Bearer "):]
		payload, err := jwt.VerifyToken(tokenString)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"Message": "Unauthorized",
				"Data":    nil,
			})
			return
		}

		ctx = SetUserID(ctx, payload.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
