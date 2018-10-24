package router

import (
	"net/http"
	//"conductor/pkg/api/application"

	log "github.com/sirupsen/logrus"
)

type Router struct{}

func (rr *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Infof("[%s] %s", r.RemoteAddr, r.URL.Path)
}
