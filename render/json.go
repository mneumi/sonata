package render

import (
	"encoding/json"
	"net/http"
)

type JSON struct {
	Data any
}

func (j *JSON) Render(w http.ResponseWriter) error {
	jsonData, err := json.Marshal(j.Data)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonData)
	return err
}

func (j *JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, "application/json; charset=utf8")
}

func (j *JSON) WriteHeader(status int, w http.ResponseWriter) {
	w.WriteHeader(status)
}
