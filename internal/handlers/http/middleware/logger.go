package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/urfave/negroni"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseWriter := negroni.NewResponseWriter(w)
		next.ServeHTTP(responseWriter, r)

		status := responseWriter.Status()

		slog.Debug(
			"http request",
			slog.Int("status", status),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("elapsed", time.Since(start).String()),
			slog.Int("size", responseWriter.Size()),
		)
	})
}
