package util

import (
	"encoding/json"
	"github.com/revel/revel"
)

type Util struct {
}

// AppResponse Represent each request response structure
type AppResponse struct {

	// status code for the particular operation
	// 200: ok, 400: not found, 500: internal server error, 401: unauthorized
	Code int

	// message to the end user for the particular operation
	Message string

	// actual result for the particular operation
	// it must be of any type
	// +optional
	Result interface{}
}

// generate generic response model for the each request
// its return response ready JSON
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
