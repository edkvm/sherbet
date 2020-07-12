package entryparser

import (
	"net/http"
	"strings"
)

func NewOpEngineMiddelware() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		rawCmd := strings.Split(r.URL.Path, "/")

		accept := r.Header.Get("Accept")

		if accept != "" && strings.Contains(accept, "image/webp") {

		}

		// TODO(edkvm): Move webP to system Env
		output(rawCmd, w)
		return
	}

	return http.HandlerFunc(fn)
}