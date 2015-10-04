package main

import (
	"flag"
	"log"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/k0kubun/pp"
	"github.com/mix3/meirin/auth"
	_ "github.com/mix3/meirin/auth/github"
	_ "github.com/mix3/meirin/auth/google"
	"github.com/mix3/meirin/config"
	"github.com/mix3/meirin/proxy"
	"github.com/mix3/ran"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "config.hcl", "config file path")
	flag.Parse()
}

func main() {
	c, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	pp.Println(c)

	cookieStore := cookiestore.New([]byte(c.SessionKey))
	cookieStore.Options(sessions.Options{Domain: c.Domain})

	a := auth.New(c)

	n := negroni.New()
	n.Use(sessions.Sessions("session", cookieStore))
	n.Use(a.NewOauth2Provider())
	n.Use(a.LoginRequired())
	n.Use(a.RestrictRequest())
	n.Use(proxy.Proxy(c))

	ran.Run(c.Addr, n)
}
