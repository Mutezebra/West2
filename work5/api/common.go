package api

import (
	"encoding/json"
	"five/pkg/ctl"
	"five/pkg/e"
)

func ErrorResponse(err error) *ctl.Response {
	_, ok := err.(*json.UnmarshalTypeError)
	if ok {
		return ctl.RespError(e.JsonUnmarshalFailed, err)
	}
	return ctl.RespError(e.ERROR, err)
}
