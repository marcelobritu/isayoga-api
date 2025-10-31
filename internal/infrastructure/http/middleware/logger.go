package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/marcelobritu/isayoga-api/pkg/logger"
	"go.uber.org/zap"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		
		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		
		logger.Info("HTTP Request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Int("status", ww.Status()),
			zap.Int("bytes", ww.BytesWritten()),
			zap.Duration("duration", duration),
			zap.String("user_agent", r.UserAgent()),
		)
	})
}
