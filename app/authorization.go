package app

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/revel/revel"
)


// filter each incoming request to the server
// also bypass some of the url/resource without auth token
func AuthorizationFilter(c *revel.Controller, fc []revel.Filter) {
	if !revel.Config.BoolDefault(AUTHENTICATED_CONF, false) {
		fc[0](c, fc[1:])
		return
	}
	if strings.Index(strings.TrimRight(c.Request.GetPath(), "/"), "/ui") == 0 {
		fc[0](c, fc[1:])
		return
	}

	permsList := c.Args[PERMS]
	if permsList == nil {
		c.Response.SetStatus(http.StatusUnauthorized)
		c.Response.GetWriter().Write([]byte("user does not have the permission to operate on this resource"))
		return
	}
	revel.AppLog.Debugf("before type conv=%+v %s", permsList, reflect.TypeOf(permsList))
	if permListInt, ok := permsList.([]interface{}); ok {
		for _, permInt := range permListInt {
			perm := map[string]string{}
			for k, v := range permInt.(map[string]interface{}) {
				if v == nil {
					continue
				}
				if b, ok := v.(bool); ok {
					if b {
						perm[k] = "true"
						continue
					}
					perm[k] = "false"
					continue
				}
				if i, ok := v.(int); ok {
					perm[k] = fmt.Sprintf("%d", i)
					continue
				}
				if f, ok := v.(float64); ok {
					perm[k] = fmt.Sprintf("%f", f)
					continue
				}
				perm[k] = v.(string)
			}
			revel.AppLog.Debugf("after type conv=%+v", perm)
			if perm["Name"] == "all" && perm["Resource"] == "*" {
				fc[0](c, fc[1:])
				return
			}

			if strings.Index(c.Request.GetPath(), "/api/v1/"+perm["Resource"]) == 0 {
				if c.Request.Method == "GET" {
					if perm["Read"] == "true" {
						fc[0](c, fc[1:])
						return
					}
				}
				if c.Request.Method == "POST" {
					if perm["Create"] == "true" {
						fc[0](c, fc[1:])
						return
					}
				}
				if c.Request.Method == "PUT" {
					if perm["Update"] == "true" {
						fc[0](c, fc[1:])
						return
					}
				}
				if c.Request.Method == "DELETE" {
					if perm["Delete"] == "true" {
						fc[0](c, fc[1:])
						return
					}
				}
			}
		}
	}
	c.Response.SetStatus(http.StatusUnauthorized)
	c.Response.GetWriter().Write([]byte("this user does not have the permission to operate on this resource"))
	return
}
