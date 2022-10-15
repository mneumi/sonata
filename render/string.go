package render

import (
	"fmt"
	"net/http"
)

type String struct {
	Format string
	Data   []any
}

func (s *String) Render(w http.ResponseWriter) error {
	if len(s.Data) > 0 {
		_, err := fmt.Fprintf(w, s.Format, s.Data...)
		return err
	}
	_, err := w.Write([]byte(s.Format))
	return err
}

func (s *String) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, "text/plain; charset=utf-8")
}

func (s *String) WriteHeader(status int, w http.ResponseWriter) {
	w.WriteHeader(status)
}
