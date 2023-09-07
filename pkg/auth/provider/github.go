package provider

import (
	"strings"

	"github.com/fermyon/spin/sdk/go/config"
	"golang.org/x/oauth2"
)

type github struct{}

func (g github) getAuthConfig() (*oauth2.Config, error) {
	clientId, err := config.Get("client_id")
	if err != nil {
		return nil, err
	}

	//for auth0 clientSecret is not required
	clientSecret, err := config.Get("client_secret")
	if err != nil {
		return nil, err
	}

	rawScopes, err := config.Get("scopes")
	if err != nil {
		rawScopes = defaultScopes
	}

	endpoint := oauth2.Endpoint{
		AuthURL:   "https://github.com/login/oauth/authorize",
		TokenURL:  "https://github.com/login/oauth/access_token",
		AuthStyle: oauth2.AuthStyleInHeader,
	}

	return &oauth2.Config{
		RedirectURL:  "/auth/success",
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       strings.Split(rawScopes, ","),
		Endpoint:     endpoint,
	}, nil
}
