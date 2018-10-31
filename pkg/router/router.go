package router

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"io/ioutil"
	"net/http"
	"strings"

	"conductor/pkg/resource"

	log "github.com/sirupsen/logrus"
)

var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

var ENV = "dev"

type Router struct{}

func (rr *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rr.HandleResources(rw, r)
}

func (rr *Router) HandleResources(rw http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"remote": r.RemoteAddr,
		"path":   r.URL.Path,
		"method": r.Method,
	}).Info("")

	//if !filters(r, true) {
	//	rw.WriteHeader(401)
	//	rw.Write([]byte("Unauthorized request!"))
	//	return
	//}

	url := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	resourceType := url[0]
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
		rslt, err := handler.Create(res)
		if err != nil {
			rw.WriteHeader(500)
			rw.Write([]byte(fmt.Sprintf("could not create resource %s: %v", resourceType, err)))
			return
		} else {
			session, _ := store.Get(r, "session-name")
			// Set some session values.
			session.Values["foo"] = "bar"
			// Save it before we write to the response/return from the handler.
			session.Save(r, rw)
			rw.WriteHeader(200)
			rw.Write([]byte(fmt.Sprintf("%+v", rslt)))
			return
		}

	case "UPDATE":
		res, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(500)
			rw.Write([]byte(fmt.Sprintf("could not create resource %s: %v", resourceType, err)))
			return
		}
		rslt, err := handler.Update(res)
		if err != nil {
			rw.WriteHeader(500)
			rw.Write([]byte(fmt.Sprintf("could not create resource %s: %v", resourceType, err)))
			return
		} else {
			session, _ := store.Get(r, "session-name")
			// Set some session values.
			session.Values["foo"] = "bar"
			// Save it before we write to the response/return from the handler.
			session.Save(r, rw)
			rw.WriteHeader(200)
			rw.Write([]byte(fmt.Sprintf("Success %s: %+v", resourceType, rslt)))
			return
		}

	case "GET":

		rslt, err := handler.Get(url[1], &resource.GetOptions{})
		if err != nil {
			rw.WriteHeader(500)
			rw.Write([]byte(fmt.Sprintf("could not create resource %s: %v", resourceType, err)))
			return
		} else {
			session, _ := store.Get(r, "session-name")
			// Set some session values.
			session.Values["foo"] = "bar"
			// Save it before we write to the response/return from the handler.
			session.Save(r, rw)
			if ENV == "dev" {
				rw.Header().Add("X-Frame-Options", "SAMEORIGIN")
				rw.Header().Add("X-XSS-Protection", "1; mode=block")
				rw.Header().Add("X-Content-Type-Options", "nosniff")
				rw.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")
				rw.Header().Add("Cache-Control", "no-cache")
			}
			rw.WriteHeader(200)
			rw.Write([]byte(fmt.Sprintf("%+v", rslt)))
			return
		}

	case "DELETE":

		rslt, err := handler.Delete(url[1], &resource.DeleteOptions{})
		if err != nil {
			rw.WriteHeader(500)
			rw.Write([]byte(fmt.Sprintf("could not create resource %s: %v", resourceType, err)))
			return
		} else {
			rw.WriteHeader(200)
			rw.Write([]byte(fmt.Sprintf("%+v", rslt)))
			return
		}
	}

}

func filters(r *http.Request, isAuthTest bool) bool {

	if isAuthTest {
		//authToken := r.Header.Get("Authorization")
	}

	return false
}
