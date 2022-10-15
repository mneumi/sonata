package render

import (
	"fmt"
	"net/http"
)

type Redirect struct {
	Status   int
	Request  *http.Request
	Location string
}

func (r *Redirect) Render(w http.ResponseWriter) error {
	// StatusMultipleChoices: http code 300
	// StatusPermanentRedirect: http code 308
	if r.Status < http.StatusMultipleChoices || r.Status > http.StatusPermanentRedirect &&
		r.Status != http.StatusCreated {
		return fmt.Errorf("cannot redirect with status code %d", r.Status)
	}
	http.Redirect(w, r.Request, r.Location, r.Status)
	return nil
}

func (r *Redirect) WriteContentType(w http.ResponseWriter) {
	// do nothing
}

func (r *Redirect) WriteHeader(status int, w http.ResponseWriter) {
	// do nothing
}
