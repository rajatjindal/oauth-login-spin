package provider

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strings"

	"golang.org/x/oauth2"
)

const defaultScopes = "openid,profile,email"

type Provider interface {
	GetAuthConfig() (*oauth2.Config, error)
	GetAuthCodeChallengeAndType() (string, string, string, error)
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

func GetPlainAuthCodeChallengeAndType() (string, string, string, error) {
	random, err := GenerateRandomString(44)
	if err != nil {
		return "", "", "", err
	}

	return random, random, "plain", nil
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func GetSha256AuthCodeChallengeAndType() (string, string, string, error) {
	verifier, err := GenerateRandomString(44)
	if err != nil {
		return "", "", "", err
	}

	hx := sha256.New()
	hx.Write([]byte(verifier))
	hashed := hx.Sum(nil)

	based := base64.StdEncoding.EncodeToString(hashed)
	based1 := strings.ReplaceAll(based, "+", "-")
	based2 := strings.ReplaceAll(based1, "/", "_")

	return verifier, regexp.MustCompile("=+$").ReplaceAllString(based2, ""), "S256", nil
}
