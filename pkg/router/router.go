package router

import (
	"net/http"
	//"conductor/pkg/api/application"

	log "github.com/sirupsen/logrus"
)

type Router struct{}

func (rr *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"remote": r.RemoteAddr,
		"path":   r.URL.Path,
		"method": r.Method,
	}).Infof("new request")
}
