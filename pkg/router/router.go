package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"conductor/pkg/resource"

	log "github.com/sirupsen/logrus"
)

type Router struct{}

func (rr *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"remote": r.RemoteAddr,
		"path":   r.URL.Path,
		"method": r.Method,
	}).Info("")

	resourceType := strings.Split(strings.Trim(r.URL.Path, "/"), "/")[0]
	handler, ok := resource.ResourcesMap[resourceType]
	if !ok {
		rw.WriteHeader(404)
		rw.Write([]byte(fmt.Sprintf("Unknown resource type %s", resourceType)))
		return
	}
	switch r.Method {
	case "POST":
		res, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(500)
			rw.Write([]byte(fmt.Sprintf("could not create resource %s: %v", resourceType, err)))
			return
		}
		_, err = handler.Create(res)
		if err != nil {
			rw.WriteHeader(500)
			rw.Write([]byte(fmt.Sprintf("could not create resource %s: %v", resourceType, err)))
			return
		}
	}
}
