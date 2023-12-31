package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/rajatjindal/oauth-login-spin/pkg/auth"
	"github.com/sirupsen/logrus"
)

func main() {}

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		u, _ := url.Parse(r.Header.Get(spinhttp.HeaderFullUrl))
		os.Setenv("spin-base-url", u.Scheme+"://"+u.Host)

		h, err := auth.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.URL.Path == "/internal/login/start" {
			logrus.Info("starting oauth function, start")
			h.Login(w, r)
			return
		}

		if r.URL.Path == "/internal/login/callback" {
			logrus.Info("starting oauth function, callback")
			h.LoginCallback(w, r)
			return
		}

		fmt.Println("unknown path ", r.URL.Path)
		http.Error(w, "not found", http.StatusNotFound)
	})
}
