package util

import (
	"github.com/revel/revel"
	"encoding/json"
)

type Util struct {
}

type AppResponse struct {
	Code int
	Message string
	Result interface{}
}


func (ar AppResponse) Apply(req *revel.Request, resp *revel.Response) {
	dataResponse, err := json.Marshal(ar)
	status := ar.Code

	if err != nil {
		resp.SetStatus(500)
		status = 500
	}

	resp.SetStatus(status)
	resp.GetWriter().Write(dataResponse)
}
