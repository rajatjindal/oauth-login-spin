package provider

import (
	"fmt"
	"os"

	"golang.org/x/oauth2"
)

const defaultScopes = "openid,profile,email"

func GetAuthConfig(p string) (*oauth2.Config, error) {
	switch p {
	case "github":
		return (github{}).getAuthConfig()
	case "auth0":
		return (auth0{}).getAuthConfig()
	}

	return nil, fmt.Errorf("unknown provider %s", p)
}

func SpinBaseURL() string {
	return os.Getenv("spin-base-url")
}
