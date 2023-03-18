package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/korylprince/httputil/jsonapi"
)

// API is the current API version
const API = "2.0"
const apiPath = "/api/" + API

// Router returns a new API router
func (s *Server) Router() http.Handler {
	r := mux.NewRouter()

	apirouter := jsonapi.New(s.output, nil, nil, nil)
	r.PathPrefix(apiPath).Handler(http.StripPrefix(apiPath, apirouter))

	apirouter.Handle("POST", "/passwd",
		s.getPassword, false)

	return r
}
