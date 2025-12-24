package logx

import (
	"log"
	"net/http"
	"time"

	"example.com/atode/backend/internal/httpx"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(p)
	w.bytes += n
	return n, err
}

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(rw, r)
		rid, _ := httpx.RequestIDFromContext(r.Context())
		log.Printf("request_id=%s method=%s path=%s status=%d bytes=%d duration=%s", rid, r.Method, r.URL.Path, rw.status, rw.bytes, time.Since(start))
	})
}
