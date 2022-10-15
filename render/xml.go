package render

import (
	"encoding/xml"
	"net/http"
)

type XML struct {
	Data any
}

func (x *XML) Render(w http.ResponseWriter) error {
	return xml.NewEncoder(w).Encode(x.Data)
}

func (x *XML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, "application/xml; charset=utf-8")
}

func (x *XML) WriteHeader(status int, w http.ResponseWriter) {
	w.WriteHeader(status)
}
