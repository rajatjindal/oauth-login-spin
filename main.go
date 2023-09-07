package main

import (
	"fmt"
	"net/http"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/rajatjindal/oauth-login-spin/pkg/auth"
	"github.com/rajatjindal/oauth-login-spin/pkg/logrus"
)

func main() {}

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("starting oauth function")
		h, err := auth.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.URL.Path == "/login/start" {
			logrus.Info("starting oauth login function, callback")
			h.Login(w, r)
			return
		}

		if r.URL.Path == "/login/callback" {
			logrus.Info("starting oauth function, callback")
			h.LoginCallback(w, r)
			return
		}

		fmt.Println("unknown path ", r.URL.Path)
	})
}
