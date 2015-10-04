package config_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mix3/meirin/config"
)

func TestParse(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	data := `
addr        = ":8080"
domain      = "example.com"
session_key = "secret123"

auth {
    service       = "google"
    client_id     = "secret client id"
    client_secret = "secret client secret"
    redirect_url  = "https://example.com/oauth2callback"
    restrictions  = [
        "yourdomain.com",
        "example@gmail.com",
    ]
}

proxies {
    hoge {
        path = "/"
        dest = "http://127.0.0.1:8081"
    }
    fuga {
        path = "/"
        dest = "http://127.0.0.1:8081"
    }
}`

	if err := ioutil.WriteFile(f.Name(), []byte(data), 0644); err != nil {
		t.Error(err)
	}

	conf, err := config.LoadConfig(f.Name())
	if err != nil {
		t.Error(err)
	}

	if conf.Addr != ":8080" {
		t.Errorf("unexpected address: %s", conf.Addr)
	}
}
