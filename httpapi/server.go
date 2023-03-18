package httpapi

import (
	"io"

	"github.com/korylprince/password-portal-server/v2/db"
)

// Server represents shared resources
type Server struct {
	sisDB  db.SISDB
	userDB db.UserDB
	output io.Writer
}

// NewServer returns a new server with the given resources
func NewServer(sisDB db.SISDB, userDB db.UserDB, output io.Writer) *Server {
	return &Server{sisDB: sisDB, userDB: userDB, output: output}
}
