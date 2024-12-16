package response

import (
	"log"
	"net/http"
	"strconv"
)

const JSONContentType = "application/json"
const CSVContentType = "text/csv"
const TextPlainContentType = "text/plain"

type HTTPResponse interface {
	WriteResponse(w http.ResponseWriter)
}

type BasicResponse struct {
	Body        []byte
	StatusCode  int
	ContentType string
}

func (b *BasicResponse) WriteResponse(w http.ResponseWriter) {
	bodyLength := len(b.Body)
	w.Header().Set("Content-Type", b.ContentType)
	if b.StatusCode != 400 {
		w.Header().Set("Content-Length", strconv.Itoa(bodyLength))
	}
	w.WriteHeader(b.StatusCode)
	if _, err := w.Write(b.Body); err != nil {
		log.Println("unable to write byte.", err)
	}
}
