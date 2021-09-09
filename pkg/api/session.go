package api

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/session"
)

const returnURLParam = "returnURL"

type loginTemplateData struct {
	URL   string
	Error string
}

func redirectToLogin(w http.ResponseWriter, returnURL string, loginError string) {
	data, _ := loginUIBox.Find("login.html")
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

func getLoginHandler(w http.ResponseWriter, r *http.Request) {
	if !config.GetInstance().HasCredentials() {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	redirectToLogin(w, r.URL.Query().Get(returnURLParam), "")
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue(returnURLParam)
	if url == "" {
		url = "/"
	}

	err := manager.GetInstance().SessionStore.Login(w, r)
	if err == session.ErrInvalidCredentials {
		// redirect back to the login page with an error
		redirectToLogin(w, url, "Username or password is invalid")
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	if err := manager.GetInstance().SessionStore.Logout(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect to the login page if credentials are required
	getLoginHandler(w, r)
}
