package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-oauth2"
	"github.com/mix3/meirin/auth"
	"github.com/mix3/meirin/config"
)

func init() {
	auth.Register("github", &GithubDriverFactory{})
}

type GithubDriverFactory struct{}

func (f GithubDriverFactory) New(c *config.Config) auth.Driver {
	return &GithubDriver{
		Config: c,
	}
}

type GithubDriver struct {
	Config *config.Config
}

func (d *GithubDriver) NewOauth2Provider() negroni.Handler {
	return oauth2.Github(&oauth2.Config{
		ClientID:     d.Config.Auth.ClientID,
		ClientSecret: d.Config.Auth.ClientSecret,
		RedirectURL:  d.Config.Auth.RedirectUrl,
		Scopes:       []string{"read:org"},
	})
}

func (d *GithubDriver) LoginRequired() negroni.Handler {
	return oauth2.LoginRequired()
}

func (d *GithubDriver) RestrictRequest() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		t := oauth2.GetToken(r)
		if d.isForbidden(t) {
			auth.Forbidden(w)
		} else {
			next(w, r)
		}
	}
}

type org []struct {
	Login string `json:"login"`
}

func (d *GithubDriver) getOrg(tokens oauth2.Tokens) org {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user/orgs", d.Config.Auth.ApiEndpoint), nil)
	if err != nil {
		log.Printf("failed to create a request to retrieve organizations: %s", err)
		return nil
	}

	req.SetBasicAuth(tokens.Access(), "x-oauth-basic")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("failed to retrieve organizations: %s", err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("failed to read body of GitHub response: %s", err)
		return nil
	}

	var result org
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("failed to decode json: %s", err)
		return nil
	}

	return result
}

func (d *GithubDriver) isForbidden(tokens oauth2.Tokens) bool {
	if len(d.Config.Auth.Restrictions) <= 0 {
		return false
	}

	info := d.getOrg(tokens)
	if info == nil {
		return true
	}

	for _, userOrg := range info {
		for _, org := range d.Config.Auth.Restrictions {
			if org == userOrg.Login {
				return false
			}
		}
	}

	log.Printf("not a member of designated organizations")
	return true
}
