package provider

import (
	"strings"

	"github.com/fermyon/spin/sdk/go/config"
	"golang.org/x/oauth2"
)

type spotify struct{}

func (g spotify) GetAuthConfig() (*oauth2.Config, error) {
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
		rawScopes = "user-read-private,user-read-email,user-top-read"
	}

	endpoint := oauth2.Endpoint{
		AuthURL:   "https://accounts.spotify.com/authorize",
		TokenURL:  "https://accounts.spotify.com/api/token",
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

func (g spotify) GetAuthCodeChallengeAndType() (string, string) {
	return GetSha256AuthCodeChallengeAndType()
}
