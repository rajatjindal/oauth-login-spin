package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/fermyon/spin/sdk/go/config"
	"github.com/google/uuid"
	"github.com/rajatjindal/oauth-login-spin/pkg/auth/provider"
	"github.com/rajatjindal/oauth-login-spin/pkg/cache"
	"github.com/rajatjindal/oauth-login-spin/pkg/cache/kvcache"
	"github.com/rajatjindal/oauth-login-spin/pkg/logrus"
	"golang.org/x/oauth2"
)

type Handler struct {
	AuthConfig     *oauth2.Config
	ChallengeCache cache.Provider

	errorURL string
}

func New() (*Handler, error) {
	authProviderName, err := config.Get("auth_provider")
	if err != nil {
		return nil, err
	}

	authConfig, err := provider.GetAuthConfig(authProviderName)
	if err != nil {
		return nil, err
	}

	errorURL, _ := config.Get("error_url")
	if errorURL == "" {
		errorURL = "/login/error"
	}

	return &Handler{
		AuthConfig:     authConfig,
		ChallengeCache: kvcache.Provider(1*time.Minute, 2*time.Minute),
		errorURL:       errorURL,
	}, nil
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	logrus.Info("starting login function")
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	oauthState := generateStateOauthCookie(state, w)
	logrus.Info("starting login function, after generating state oauth cookie")

	/*
		AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
		validate that it matches the the state query parameter on your redirect callback.
	*/

	logrus.Info("starting login function, creating challenge code")
	challengeCode := uuid.New().String() + uuid.New().String()
	logrus.Info("login function, storing challenge in cache ", state, challengeCode)
	err := h.storeChallenge(state, challengeCode)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	logrus.Info("starting login function, caching challenge code")

	params := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("code_challenge", challengeCode),
		oauth2.SetAuthURLParam("code_challenge_method", "plain"),
	}

	logrus.Info("starting login function, redirecting to auth code url")
	u := h.AuthConfig.AuthCodeURL(oauthState, params...)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (h *Handler) LoginCallback(w http.ResponseWriter, r *http.Request) {
	logrus.Info("starting login callback function")
	oauthState, err := r.Cookie("oauthstate")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Info("login callback function, after getting state from cookie")

	if r.FormValue("state") != oauthState.Value {
		logrus.Infof("invalid oauth state %s", r.FormValue("state"))
		http.Redirect(w, r, h.errorURL, http.StatusTemporaryRedirect)
		return
	}

	logrus.Info("login callback function, after getting state from form value ", oauthState.Value)
	challengeCode := h.getChallenge(oauthState.Value)

	logrus.Info("login callback function, after getting challenge code from auth state ", oauthState.Value, challengeCode)

	token, err := h.exchangeToken(challengeCode, r.FormValue("code"))
	if err != nil {
		logrus.Info("login callback function, error calling exchange token", err.Error())
		http.Redirect(w, r, h.errorURL, http.StatusTemporaryRedirect)
		return
	}

	logrus.Info("login callback function, doing redirect now")
	link := fmt.Sprintf("%s/#access-token=%s", h.AuthConfig.RedirectURL, token)
	http.Redirect(w, r, link, http.StatusTemporaryRedirect)
}

func generateStateOauthCookie(state string, w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)
	return state
}

func (h *Handler) exchangeToken(challenge, code string) (string, error) {
	logrus.Info("exchange token function")
	params := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("code_verifier", challenge),
	}

	logrus.Info("exchange token function, starting the exchange")
	token, err := h.AuthConfig.Exchange(context.Background(), code, params...)
	if err != nil {
		return "", fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	logrus.Infof("exchange token function, exchange done. token is %#v", token)

	return token.Extra("access_token").(string), nil
}

func (h *Handler) storeChallenge(state, challenge string) error {
	return h.ChallengeCache.Set(state, challenge)
}

func (h *Handler) getChallenge(state string) string {
	v, ok := h.ChallengeCache.Get(state)
	if !ok {
		return ""
	}

	return v.(string)
}
