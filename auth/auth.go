package auth

import (
	"log"
	"net/http"
	"sync"

	"github.com/codegangsta/negroni"
	"github.com/mix3/meirin/config"
)

var (
	factoriesMu sync.Mutex
	factories   = make(map[string]DriverFactory)
)

type DriverFactory interface {
	New(c *config.Config) Driver
}

type Driver interface {
	NewOauth2Provider() negroni.Handler
	LoginRequired() negroni.Handler
	RestrictRequest() negroni.HandlerFunc
}

func Register(name string, factory DriverFactory) {
	factoriesMu.Lock()
	defer factoriesMu.Unlock()

	if factory == nil {
		log.Fatal("auth: Register factory is nil")
	}

	if _, dup := factories[name]; dup {
		log.Fatal("auth: Register called twice for factory " + name)
	}

	factories[name] = factory
}

func Forbidden(w http.ResponseWriter) {
	w.WriteHeader(403)
	w.Write([]byte("Access denied"))
}

func New(c *config.Config) Driver {
	d, ok := factories[c.Auth.Service]
	if !ok {
		log.Fatal("auth: unknown driver %s (forgotten import?)", c.Auth.Service)
	}
	return d.New(c)
}
