package provider

import (
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

const defaultScopes = "openid,profile,email"

type Provider interface {
	GetAuthConfig() (*oauth2.Config, error)
	GetAuthCodeChallengeAndType() (string, string)
}

func GetProvider(p string) (Provider, error) {
	switch p {
	case "github":
		return &github{}, nil
	case "auth0":
		return &auth0{}, nil
	case "spotify":
		return &spotify{}, nil
	}

	return nil, fmt.Errorf("unknown provider %s", p)
}

func SpinBaseURL() string {
	return os.Getenv("spin-base-url")
}

func GetPlainAuthCodeChallengeAndType() (string, string) {
	return uuid.New().String() + uuid.New().String(), "plain"
}

func GetSha256AuthCodeChallengeAndType() (string, string) {
	challengeCode := uuid.New().String() + uuid.New().String()

	hx := sha256.New()
	hx.Write([]byte(challengeCode))
	bs := hx.Sum(nil)

	return fmt.Sprintf("%x\n", bs), "S256"
}
