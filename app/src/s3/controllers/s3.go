package controllers

import (
	"net/url"

	static "github.com/revel/modules/static/app/controllers"
	"github.com/revel/revel"
)

type S3 struct {
	static.Static
}

// Serve the resource of the UI from local or AWS-S3 based on user's selection
func (s S3) ServeUrlOrLocal(filepath string) revel.Result {
	if filepath == "" {
		filepath = "index.html"
	}
	localPath := revel.Config.BoolDefault("koki.ui.local", false)
	if localPath {
		s.Params = new(revel.Params)
		s.Params.Fixed = url.Values{}
		s.Params.Fixed.Set("prefix", "dist")
		return s.Serve("dist", filepath)
	}
	redirectURL := revel.Config.StringDefault("koki.ui.url", "http://ui.koki.io.s3-website-us-east-1.amazonaws.com/")
	if last := redirectURL[len(redirectURL)-1:]; last != "/" {
		redirectURL = redirectURL + "/"
	}
	return s.Redirect(redirectURL + filepath)
}
