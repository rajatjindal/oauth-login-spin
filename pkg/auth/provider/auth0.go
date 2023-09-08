package provider

import (
	"fmt"
	"strings"

	"github.com/fermyon/spin/sdk/go/config"
	"golang.org/x/oauth2"
)

type auth0 struct{}

func (a auth0) getAuthConfig() (*oauth2.Config, error) {
	tenant, err := config.Get("tenant")
	if err != nil {
		return nil, err
	}

	clientId, err := config.Get("client_id")
	if err != nil {
		return nil, err
	}

	//for auth0 clientSecret is not required
	clientSecret, _ := config.Get("client_secret")

	rawScopes, err := config.Get("scopes")
	if err != nil {
		rawScopes = defaultScopes
	}

	endpoint := oauth2.Endpoint{
		AuthURL:   fmt.Sprintf("https://%s.auth0.com/authorize", tenant),
		TokenURL:  fmt.Sprintf("https://%s.auth0.com/oauth/token", tenant),
		AuthStyle: oauth2.AuthStyleInHeader,
	}

	return &oauth2.Config{
		RedirectURL:  SpinBaseURL() + "/internal/login/callback",
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       strings.Split(rawScopes, ","),
		Endpoint:     endpoint,
	}, nil
}
