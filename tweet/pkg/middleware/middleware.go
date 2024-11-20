package middleware

import (
	"fmt"
	"net/http"

	"github.com/jackgris/twitter-backend/tweet/pkg/logger"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// LogResponse is the middleware layer to log all the HTTP requests
func LogResponse(next http.HandlerFunc, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, req)

		statusCode := lrw.statusCode
		log.Info(req.Context(), fmt.Sprintf("--> %s %s  Status %d %s", req.Method, req.URL.Path, statusCode, http.StatusText(statusCode)))
	}
}
