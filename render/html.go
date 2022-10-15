package render

import (
	"net/http"
	"text/template"
)

type HTML struct {
	Data       any
	Name       string
	Template   *template.Template
	IsTemplate bool
}

type HTMLRender struct {
	Template *template.Template
}

func (h *HTML) Render(w http.ResponseWriter) error {
	if h.IsTemplate {
		err := h.Template.ExecuteTemplate(w, h.Name, h.Data)
		return err
	}
	_, err := w.Write([]byte(h.Data.(string)))
	return err
}

func (h *HTML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, "text/html; charset=utf8")
}

func (h *HTML) WriteHeader(status int, w http.ResponseWriter) {
	w.WriteHeader(status)
}
