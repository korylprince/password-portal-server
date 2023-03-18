package ldap

import (
	"fmt"
	"log"
	"strings"

	adauth "github.com/korylprince/go-ad-auth/v3"
	"github.com/korylprince/password-portal-server/v2/db"
	"github.com/korylprince/securetoken"
)

// DB represents a connection to an Active Directory server
type DB struct {
	config   *adauth.Config
	bindUser string
	bindPass string
	key      []byte
	debug    bool
}

// New returns a new *DB with the given parameters
func New(config *adauth.Config, username, password, key string, debug bool) *DB {
	return &DB{config: config, bindUser: username, bindPass: password, key: []byte(key), debug: debug}
}

// Bind returns a bound connection to an Active Directory server
func (d *DB) Bind() (*adauth.Conn, error) {
	conn, err := d.config.Connect()
	if err != nil {
		return nil, err
	}

	status, err := conn.Bind(d.bindUser, d.bindPass)
	if err != nil {
		return nil, err
	}

	if !status {
		return nil, fmt.Errorf("Invalid bind credentials for user %s", d.bindUser)
	}

	return conn, nil
}

// Get returns the user with the given id, nil if the user doesn't exist, or an error if one occurred
func (d *DB) Get(id string) (*db.User, error) {
	conn, err := d.Bind()
	if err != nil {
		return nil, fmt.Errorf("Error binding to server: %v", err)
	}
	defer conn.Conn.Close()

	entry, err := conn.GetAttributes("employeeID", id, []string{"sn", "givenname", "sAMAccountName", "adminDescription"})
	if err != nil {
		if strings.HasSuffix(err.Error(), "no entries returned") {
			return nil, nil
		}
		return nil, fmt.Errorf("Error searching for user: %v", err)
	}

	if entry == nil {
		return nil, nil
	}

	pass, err := securetoken.DecryptToken(entry.GetRawAttributeValue("adminDescription"), d.key, 0)
	if err != nil {
		pass = []byte("")
		if d.debug {
			log.Printf("Unable to decrypt password for user %s: %v\n", entry.GetAttributeValue("sAMAccountName"), err)
		}
	}

	user := &db.User{
		FirstName: entry.GetAttributeValue("givenName"),
		LastName:  entry.GetAttributeValue("sn"),
		Username:  entry.GetAttributeValue("sAMAccountName"),
		Password:  string(pass),
	}

	return user, nil
}
