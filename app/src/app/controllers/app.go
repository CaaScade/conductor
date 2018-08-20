package controllers

import (
	"io/ioutil"
	"net/http"
	"net/url"

	static "github.com/revel/modules/static/app/controllers"
	"github.com/revel/revel"
)

type MyHtml string

func (r MyHtml) Apply(req *revel.Request, resp *revel.Response) {
	resp.WriteHeader(http.StatusOK, "text/html")
	resp.GetWriter().Write([]byte(r))
}

type App struct {
	static.Static
}

func (c App) Login() revel.Result {
	localPath := revel.Config.BoolDefault("koki.ui.local", false)
	if localPath {
		c.Params = new(revel.Params)
		c.Params.Fixed = url.Values{}
		c.Params.Fixed.Set("prefix", "dist")
		return c.Serve("dist", "index.html")
	}

	redirectURL := revel.Config.StringDefault("koki.ui.url", "http://ui.koki.io.s3-website-us-east-1.amazonaws.com/")
	if last := redirectURL[len(redirectURL)-1:]; last != "/" {
		redirectURL = redirectURL + "/"
	}
	resp, err := http.Get(redirectURL + "/#/login")
	if err != nil {
		return revel.PlaintextErrorResult{err}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return revel.PlaintextErrorResult{err}
	}
	return MyHtml(body)
}

func (c App) Index() revel.Result {
	localPath := revel.Config.BoolDefault("koki.ui.local", false)
	if localPath {
		c.Params = new(revel.Params)
		c.Params.Fixed = url.Values{}
		c.Params.Fixed.Set("prefix", "dist")
		return c.Serve("dist", "index.html")
	}

	redirectURL := revel.Config.StringDefault("koki.ui.url", "http://ui.koki.io.s3-website-us-east-1.amazonaws.com/")
	if last := redirectURL[len(redirectURL)-1:]; last != "/" {
		redirectURL = redirectURL + "/"
	}

	resp, err := http.Get(redirectURL + "index.html")
	if err != nil {
		return revel.PlaintextErrorResult{err}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return revel.PlaintextErrorResult{err}
	}
	return MyHtml(body)
}
