package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/alexbrainman/odbc"
	"github.com/didip/tollbooth/v5"
	"github.com/gorilla/handlers"
	"github.com/kelseyhightower/envconfig"
	auth "github.com/korylprince/go-ad-auth/v3"
	"github.com/korylprince/password-portal-server/v2/db/ldap"
	"github.com/korylprince/password-portal-server/v2/db/skyward"
	"github.com/korylprince/password-portal-server/v2/httpapi"
)

func main() {
	var config = &Config{}
	err := envconfig.Process("PORTAL", config)
	if err != nil {
		log.Fatalln("Error reading configuration from environment:", err)
	}

	adConfig := &auth.Config{
		Server:   config.LDAPServer,
		Port:     config.LDAPPort,
		BaseDN:   config.LDAPBaseDN,
		Security: config.SecurityType(),
	}

	sisDB, err := skyward.New(config.SQLDSN)
	if err != nil {
		log.Fatalln("Unable to open database:", err)
	}
	userDB := ldap.New(adConfig, config.LDAPBindUPN, config.LDAPBindPassword, config.SecureTokenKey, config.Debug)

	s := httpapi.NewServer(sisDB, userDB, os.Stdout)

	rps := float64(config.LimiterCount) / (float64(config.LimiterDuration) * 60)
	limiter := tollbooth.NewLimiter(rps, nil)
	limiter.SetBurst(config.LimiterCount)

	chain := handlers.CompressHandler(handlers.ProxyHeaders(handlers.CombinedLoggingHandler(os.Stdout, http.StripPrefix(config.Prefix, tollbooth.LimitHandler(limiter, s.Router())))))

	log.Println("Listening on:", config.ListenAddr)
	log.Println(http.ListenAndServe(config.ListenAddr, chain))
}
