package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zlog.Info("handling HTTP request",
			zap.String("method", r.Method),
			zap.Any("host", r.Host),
			zap.Any("url", r.URL),
			zap.Any("headers", r.Header),
			zap.Any("body", json.NewDecoder(r.Body)),
		)
		next.ServeHTTP(w, r)
	})
}
