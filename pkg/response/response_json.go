package response

import (
	"encoding/json"
	"net/http"

	http_error "github.com/okiww/billing-loan-system/pkg/errors"
	"github.com/pkg/errors"
)

type JSONResponse struct {
	BasicResponse
	JSONBody JSONBody
	Error    error
}
type JSONBody struct {
	Message string            `json:"message,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
	Meta    interface{}       `json:"meta,omitempty"`
	Error   *http_error.Error `json:"error,omitempty"`
}

func NewJSONResponse() *JSONResponse {
	return &JSONResponse{
		BasicResponse: BasicResponse{
			ContentType: JSONContentType,
		},
	}
}

func (r *JSONResponse) SetData(data interface{}) *JSONResponse {
	r.StatusCode = http.StatusOK
	r.JSONBody.Data = data
	return r
}

func (r *JSONResponse) SetMeta(meta interface{}) *JSONResponse {
	r.JSONBody.Meta = meta
	return r
}

func (r *JSONResponse) SetMessage(message string) *JSONResponse {
	r.JSONBody.Message = message
	return r
}

func (r *JSONResponse) SetError(err error) *JSONResponse {
	respErr := &http_error.Error{}
	if errors.As(err, &respErr) {
		r.JSONBody.Error = respErr
	} else {
		// when unspecified error is provided it will categorize the response as internal server error
		r.JSONBody.Error = http_error.NewError(err.Error(), http.StatusInternalServerError)
	}
	return r
}

func (r *JSONResponse) WriteResponse(w http.ResponseWriter) {
	b, err := json.Marshal(r.JSONBody)
	if err != nil {
		JSONBody := JSONBody{
			Error: http_error.NewError(err.Error(), http.StatusInternalServerError),
		}
		b, _ = json.Marshal(JSONBody)
	}
	r.Body = b
	if r.JSONBody.Error != nil {
		r.StatusCode = r.JSONBody.Error.ErrorCode
	}
	r.BasicResponse.WriteResponse(w)
}
