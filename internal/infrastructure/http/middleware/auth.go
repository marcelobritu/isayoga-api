package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/marcelobritu/isayoga-api/pkg/auth"
	"github.com/marcelobritu/isayoga-api/pkg/logger"
	"go.uber.org/zap"
)

type contextKey string

const UserClaimsKey contextKey = "user_claims"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token de autenticação não fornecido", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Formato de token inválido", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		claims, err := auth.ValidateToken(token)
		if err != nil {
			logger.Warn("Token inválido", zap.Error(err))
			http.Error(w, "Token inválido ou expirado", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(UserClaimsKey).(*auth.Claims)
		if !ok {
			http.Error(w, "Não autorizado", http.StatusUnauthorized)
			return
		}

		if claims.Role != "admin" && claims.Role != "instructor" {
			logger.Warn("Tentativa de acesso não autorizado à área administrativa",
				zap.String("user_id", claims.UserID),
				zap.String("role", string(claims.Role)),
			)
			http.Error(w, "Acesso negado: apenas administradores e instrutores", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

