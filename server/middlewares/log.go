package middlewares

import (
	"log"
	"net/http"
	"time"
)

type HijackResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *HijackResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &HijackResponseWriter{
			ResponseWriter: w,
			statusCode:     200,
		}
		t := time.Now()
		next.ServeHTTP(rw, r)
		log.Printf("Request: %s [%d] %s (%s)\n", r.Method, rw.statusCode, r.RequestURI, time.Since(t))
	})
}
