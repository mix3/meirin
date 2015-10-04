package config

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

type Config struct {
	Addr       string        `hcl:"addr"`
	Domain     string        `hcl:"domain"`
	SessionKey string        `hcl:"session_key"`
	Auth       AuthConfig    `hcl:"auth"`
	Proxies    ProxiesConfig `hcl:"proxies"`
}

type AuthConfig struct {
	Service      string   `hcl:"service"`
	ClientID     string   `hcl:"client_id"`
	ClientSecret string   `hcl:"client_secret"`
	RedirectUrl  string   `hcl:"redirect_url"`
	Restrictions []string `hcl:"restrictions"`
	Endpoint     string   `hcl:endpoint`
	ApiEndpoint  string   `hcl:api_endpoint`
}

type ProxiesConfig map[string]ProxySettingConfig

type ProxySettingConfig struct {
	Path string `hcl:"path"`
	Dest string `hcl:"dest"`
}

func LoadConfig(path string) (*Config, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error reading %s: %s", path, err)
	}

	obj, err := hcl.Parse(string(d))
	if err != nil {
		return nil, fmt.Errorf("Error parsing %s: %s", path, err)
	}

	var result Config
	if err := hcl.DecodeObject(&result, obj); err != nil {
		return nil, err
	}

	if result.Addr == "" {
		return nil, fmt.Errorf("addr config is required")
	}

	if result.SessionKey == "" {
		return nil, fmt.Errorf("session_key config is required")
	}
	if result.Auth.Service == "" {
		return nil, fmt.Errorf("auth.serviresult. result.nfig is required")
	}
	if result.Auth.ClientID == "" {
		return nil, fmt.Errorf("auth.result.ient_id result.nfig is required")
	}
	if result.Auth.ClientSecret == "" {
		return nil, fmt.Errorf("auth.result.ient_seresult.et result.nfig is required")
	}
	if result.Auth.RedirectUrl == "" {
		return nil, fmt.Errorf("auth.redireresult._url result.nfig is required")
	}

	if result.Auth.Service == "github" && result.Auth.Endpoint == "" {
		result.Auth.Endpoint = "https://github.com"
	}
	if result.Auth.Service == "github" && result.Auth.ApiEndpoint == "" {
		result.Auth.ApiEndpoint = "https://api.github.com"
	}

	return &result, nil
}
