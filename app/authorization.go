package app

import (
	"net/http"
	"strings"

	"github.com/koki/conductor/app/models"
	"github.com/revel/revel"
)

func AuthorizationFilter(c *revel.Controller, fc []revel.Filter) {
	if !revel.Config.BoolDefault(AUTHENTICATED_CONF, false) {
		fc[0](c, fc[1:])
		return
	}
	if strings.Index(strings.TrimRight(c.Request.GetPath(), "/"), "/ui") == 0 {
		fc[0](c, fc[1:])
		return
	}

	perms := c.Args[PERMS]
	if perms == nil {
		c.Response.SetStatus(http.StatusUnauthorized)
		c.Response.GetWriter().Write([]byte("user does not have the permission to operate on this resource"))
		return
	}
	if _, ok := perms.([]models.Permission); ok {
		for _, perm := range c.Args[PERMS].([]models.Permission) {
			if perm.Resource == "all" {
				fc[0](c, fc[1:])
				return
			}

			if strings.Index(c.Request.GetPath(), "/api/v1/"+perm.Resource) == 0 {
				if c.Request.Method == "GET" {
					if perm.Read {
						fc[0](c, fc[1:])
						return
					}
				}
				if c.Request.Method == "POST" {
					if perm.Create {
						fc[0](c, fc[1:])
						return
					}
				}
				if c.Request.Method == "PUT" {
					if perm.Update {
						fc[0](c, fc[1:])
						return
					}
				}
				if c.Request.Method == "DELETE" {
					if perm.Delete {
						fc[0](c, fc[1:])
						return
					}
				}
			}
		}
	}
	c.Response.SetStatus(http.StatusUnauthorized)
	c.Response.GetWriter().Write([]byte("user does not have the permission to operate on this resource"))
	return
}
