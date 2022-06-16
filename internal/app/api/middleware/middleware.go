package middleware

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		//log.Println(r.RequestURI)
		t1 := time.Now()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.WithFields(log.Fields{
			"url":         r.URL.String(),
			"user-agent":  r.Header["User-Agent"],
			"method":      r.Method,
			"remote-addr": r.RemoteAddr,
			"duration":    t2.Sub(t1).String(),
		}).Info()
	})
}
