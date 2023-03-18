package httpapi

import (
	"net/http"
	"time"

	"github.com/korylprince/httputil/jsonapi"
)

func (s *Server) getPassword(r *http.Request) (int, interface{}) {
	type request struct {
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		BirthDate time.Time `json:"birth_date"`
		SSN       string    `json:"ssn"`
	}

	req := new(request)
	if err := jsonapi.ParseJSONBody(r, req); err != nil {
		return http.StatusBadRequest, err
	}

	id, err := s.sisDB.GetID(req.FirstName, req.LastName, req.SSN, req.BirthDate)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	user, err := s.userDB.Get("s" + id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if user == nil {
		return http.StatusNotFound, nil
	}

	jsonapi.LogActionID(r, "s"+id)

	return http.StatusOK, user
}
