package httperror

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

type Responder interface {
	Response() (int, optimus-ide-collabsdk.Response)
}

func IsResponder(err error) (Responder, bool) {
	var responseErr Responder
	if errors.As(err, &responseErr) {
		return responseErr, true
	}
	return nil, false
}

func NewResponseError(status int, resp optimus-ide-collabsdk.Response) error {
	return &responseError{
		status:   status,
		response: resp,
	}
}

func WriteResponseError(ctx context.Context, rw http.ResponseWriter, err error) {
	if responseErr, ok := IsResponder(err); ok {
		code, resp := responseErr.Response()

		httpapi.Write(ctx, rw, code, resp)
		return
	}

	httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
		Message: "Internal server error",
		Detail:  err.Error(),
	})
}

type responseError struct {
	status   int
	response optimus-ide-collabsdk.Response
}

var (
	_ error     = (*responseError)(nil)
	_ Responder = (*responseError)(nil)
)

func (e *responseError) Error() string {
	return fmt.Sprintf("%s: %s", e.response.Message, e.response.Detail)
}

func (e *responseError) Status() int {
	return e.status
}

func (e *responseError) Response() (int, optimus-ide-collabsdk.Response) {
	return e.status, e.response
}

var ErrResourceNotFound = NewResponseError(http.StatusNotFound, httpapi.ResourceNotFoundResponse)
