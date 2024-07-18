package middleware

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func WithHTTPLoggingMiddleware(handler http.Handler) http.Handler {
	fn := func(resp http.ResponseWriter, req *http.Request) {
		log.Printf("[middleware.Logging] start: %v", req.Method)
		handler.ServeHTTP(resp, req)
		log.Printf("[middleware.Logging] end: %v", req.ContentLength)
	}
	return http.HandlerFunc(fn)

}