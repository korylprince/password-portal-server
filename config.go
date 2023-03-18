package main

import (
	"log"
	"strings"

	auth "github.com/korylprince/go-ad-auth/v3"
)

// Config represents options given in the environment
type Config struct {
	LDAPServer   string `required:"true"`
	LDAPPort     int    `required:"true" default:"389"`
	LDAPBaseDN   string `required:"true"`
	LDAPSecurity string `required:"true" default:"none"`

	LDAPBindUPN      string `required:"true"`
	LDAPBindPassword string `required:"true"`

	SQLDSN string `required:"true"`

	LimiterCount    int `required:"true" default:"10"`
	LimiterDuration int `required:"true" default:"5"` //in minutes

	SecureTokenKey string `required:"true"`

	ListenAddr string `required:"true" default:":80"` //addr format used for net.Dial; required
	Prefix     string //url prefix to mount api to without trailing slash
	Debug      bool   `default:"false"`
}

// SecurityType returns the auth.SecurityType for the config
func (c *Config) SecurityType() auth.SecurityType {
	switch strings.ToLower(c.LDAPSecurity) {
	case "", "none":
		return auth.SecurityNone
	case "tls":
		return auth.SecurityTLS
	case "starttls":
		return auth.SecurityStartTLS
	default:
		log.Fatalln("Invalid PORTAL_LDAPSECURITY:", c.LDAPSecurity)
	}
	panic("unreachable")
}
