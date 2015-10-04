package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/codegangsta/negroni"
	"github.com/mix3/meirin/config"
)

func Proxy(c *config.Config) negroni.HandlerFunc {
	proxyMap := make(map[string]*url.URL)
	for k, v := range c.Proxies {
		u, err := url.Parse(fmt.Sprintf(v.Dest))
		if err != nil {
			log.Fatal(err)
		}

		proxyMap[fmt.Sprintf("%s.%s", k, c.Domain)] = u
	}

	return func(w http.ResponseWriter, r *http.Request, _ http.HandlerFunc) {
		dest, ok := proxyMap[r.Host]
		if !ok {
			http.NotFound(w, r)
			return
		}

		httputil.NewSingleHostReverseProxy(dest).ServeHTTP(w, r)
	}
}
