package app

import (
	"github.com/gin-gonic/gin"
)

// GoatResponse defines the structure of response data
type GoatResponse struct {
	Context *gin.Context `json:"-"`
	Code    int          `json:"code"`
	Msg     string       `json:"msg"`
	Data    interface{}  `json:"data"`
}

const (
	SUCCESS        = 200
	CREATED        = 201
	ERROR          = 500
	INVALID_PARAMS = 400
	NOT_AUTHORIED  = 403
	NOT_FOUND      = 404
)

var Msg = map[int]string{
	SUCCESS:        "OK",
	CREATED:        "Created",
	ERROR:          "Fail",
	INVALID_PARAMS: "Invalid Parameters",
	NOT_AUTHORIED:  "Client Not Authorized",
	NOT_FOUND:      "Resource Not Found",
}

// GetMsg returns the message to corresponding error code
func GetMsg(code int) (m string) {
	m, ok := Msg[code]
	if !ok {
		return Msg[ERROR]
	}
	return m
}

func (c GoatResponse) Response(status_code, err_code int, data interface{}) {
	c.Context.JSON(status_code, GoatResponse{
		Code: err_code,
		Msg:  GetMsg(err_code),
		Data: data,
	})
}
