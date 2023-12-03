package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/fermyon/spin/sdk/go/config"
	"github.com/rajatjindal/oauth-login-spin/pkg/auth/provider"
	"github.com/rajatjindal/oauth-login-spin/pkg/cache"
	"github.com/rajatjindal/oauth-login-spin/pkg/cache/kvcache"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Handler struct {
	ChallengeCache cache.Provider

	authprovider provider.Provider
	authconfig   *oauth2.Config
	successURL   string
	errorURL     string
}

func New() (*Handler, error) {
	authProviderName, err := config.Get("auth_provider")
	if err != nil {
		return nil, err
	}

	authprovider, err := provider.GetProvider(authProviderName)
	if err != nil {
		return nil, err
	}

	authconfig, err := authprovider.GetAuthConfig()
	if err != nil {
		return nil, err
	}

	errorURL, err := config.Get("error_url")
	if err != nil {
		return nil, err
	}

	successURL, err := config.Get("success_url")
	if err != nil {
		return nil, err
	}

	return &Handler{
		ChallengeCache: kvcache.Provider(1*time.Minute, 2*time.Minute),

		authprovider: authprovider,
		authconfig:   authconfig,
		successURL:   successURL,
		errorURL:     errorURL,
	}, nil
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	oauthState := generateStateOauthCookie(state, w)

	/*
		AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
		validate that it matches the the state query parameter on your redirect callback.
	*/

	verifier, challengeCode, challengeType, err := h.authprovider.GetAuthCodeChallengeAndType()
	if err != nil {
		logrus.Error("login start function, failed to store challenge code in cache")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if challengeType == "S256" {
		err = h.storeChallenge(state, verifier)
	} else {
		err = h.storeChallenge(state, challengeCode)
	}

	if err != nil {
		logrus.Error("login start function, failed to store challenge code in cache")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	params := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("code_challenge", challengeCode),
		oauth2.SetAuthURLParam("code_challenge_method", challengeType),
	}

	u := h.authconfig.AuthCodeURL(oauthState, params...)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (h *Handler) LoginCallback(w http.ResponseWriter, r *http.Request) {
	oauthState, err := r.Cookie("oauthstate")
	if err != nil {
		logrus.Errorf("login callback function, failed to get oauthstate cookie. error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.FormValue("state") != oauthState.Value {
		logrus.Errorf("login callback function, invalid oauth state %s", r.FormValue("state"))
		http.Redirect(w, r, h.errorURL, http.StatusTemporaryRedirect)
		return
	}

	challengeCode := h.getChallenge(oauthState.Value)
	token, err := h.exchangeToken(challengeCode, r.FormValue("code"))
	if err != nil {
		logrus.Errorf("login callback function, error calling exchange token. error: %v", err.Error())
		http.Redirect(w, r, h.errorURL, http.StatusTemporaryRedirect)
		return
	}

	logrus.Info("login successful, redirecting now")
	link := fmt.Sprintf("%s#access-token=%s", h.successURL, token)
	http.Redirect(w, r, link, http.StatusTemporaryRedirect)
}

func generateStateOauthCookie(state string, w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)
	return state
}

func (h *Handler) exchangeToken(challenge, code string) (string, error) {
	params := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("code_verifier", challenge),
	}

	token, err := h.authconfig.Exchange(context.Background(), code, params...)
	if err != nil {
		return "", fmt.Errorf("code exchange wrong: %s", err.Error())
	}

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
