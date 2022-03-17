package api

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/session"
)

const loginRootDir = "login"
const returnURLParam = "returnURL"

func getLoginPage(loginUIBox embed.FS) []byte {
	data, err := loginUIBox.ReadFile(loginRootDir + "/login.html")
	if err != nil {
		panic(err)
	}
	return data
}

type loginTemplateData struct {
	URL   string
	Error string
}

func redirectToLogin(loginUIBox embed.FS, w http.ResponseWriter, returnURL string, loginError string) {
	data := getLoginPage(loginUIBox)
	templ, err := template.New("Login").Parse(string(data))
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err), http.StatusInternalServerError)
		return
	}

	err = templ.Execute(w, loginTemplateData{URL: returnURL, Error: loginError})
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err), http.StatusInternalServerError)
	}
}

func getLoginHandler(loginUIBox embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !config.GetInstance().HasCredentials() {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		redirectToLogin(loginUIBox, w, r.URL.Query().Get(returnURLParam), "")
	}
}

func handleLogin(loginUIBox embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.FormValue(returnURLParam)
		if url == "" {
			url = "/"
		}

		err := manager.GetInstance().SessionStore.Login(w, r)
		if errors.Is(err, session.ErrInvalidCredentials) {
			// redirect back to the login page with an error
			redirectToLogin(loginUIBox, w, url, "Username or password is invalid")
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}

func handleLogout(loginUIBox embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := manager.GetInstance().SessionStore.Logout(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// redirect to the login page if credentials are required
		getLoginHandler(loginUIBox)(w, r)
	}
}
