package google

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-oauth2"
	"github.com/mix3/meirin/auth"
	"github.com/mix3/meirin/config"
)

func init() {
	auth.Register("google", &GoogleDriverFactory{})
}

type GoogleDriverFactory struct{}

func (f GoogleDriverFactory) New(c *config.Config) auth.Driver {
	return &GoogleDriver{
		Config: c,
	}
}

type GoogleDriver struct {
	Config *config.Config
}

func (d *GoogleDriver) NewOauth2Provider() negroni.Handler {
	return oauth2.Google(&oauth2.Config{
		ClientID:     d.Config.Auth.ClientID,
		ClientSecret: d.Config.Auth.ClientSecret,
		RedirectURL:  d.Config.Auth.RedirectUrl,
		Scopes:       []string{"email"},
	})
}

func (d *GoogleDriver) LoginRequired() negroni.Handler {
	return oauth2.LoginRequired()
}

func (d *GoogleDriver) RestrictRequest() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		t := oauth2.GetToken(r)
		if d.isForbidden(t) {
			auth.Forbidden(w)
		} else {
			next(w, r)
		}
	}
}

type tokeninfo struct {
	Email string `json:"email"`
}

func (d *GoogleDriver) getTokeninfo(tokens oauth2.Tokens) *tokeninfo {
	res, err := http.Get("https://www.googleapis.com/oauth2/v1/tokeninfo?access_token=" + tokens.Access())
	defer res.Body.Close()
	if err != nil {
		return nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}

	var result tokeninfo
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil
	}

	return &result
}

func (d *GoogleDriver) isForbidden(tokens oauth2.Tokens) bool {
	if len(d.Config.Auth.Restrictions) <= 0 {
		return false
	}

	info := d.getTokeninfo(tokens)
	if info == nil || info.Email == "" {
		log.Printf("email not found")
		return true
	}

	for _, d := range d.Config.Auth.Restrictions {
		if strings.Contains(d, "@") && d == info.Email {
			log.Printf("user %s logged in", info.Email)
			return false
		}
		if strings.HasSuffix(info.Email, "@"+d) {
			log.Printf("user %s logged in", info.Email)
			return false
		}
	}

	log.Printf("email doesn't allow: %s", info.Email)
	return true
}
