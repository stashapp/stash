package api

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/session"
	"github.com/stashapp/stash/pkg/utils"
	"github.com/stashapp/stash/ui"
)

const returnURLParam = "returnURL"

func getLoginPage() []byte {
	data, err := fs.ReadFile(ui.LoginUIBox, "login.html")
	if err != nil {
		panic(err)
	}
	return data
}

type loginTemplateData struct {
	URL   string
	Error string
}

func serveLoginPage(w http.ResponseWriter, r *http.Request, returnURL string, loginError string) {
	loginPage := string(getLoginPage())
	prefix := getProxyPrefix(r)
	loginPage = strings.ReplaceAll(loginPage, "/%BASE_URL%", prefix)

	templ, err := template.New("Login").Parse(loginPage)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err), http.StatusInternalServerError)
		return
	}

	buffer := bytes.Buffer{}
	err = templ.Execute(&buffer, loginTemplateData{URL: returnURL, Error: loginError})
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	// we shouldn't need to set plugin exceptions here
	setPageSecurityHeaders(w, r, nil)

	utils.ServeStaticContent(w, r, buffer.Bytes())
}

func handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		returnURL := r.URL.Query().Get(returnURLParam)

		if !config.GetInstance().HasCredentials() {
			if returnURL != "" {
				http.Redirect(w, r, returnURL, http.StatusFound)
			} else {
				prefix := getProxyPrefix(r)
				http.Redirect(w, r, prefix+"/", http.StatusFound)
			}
			return
		}

		serveLoginPage(w, r, returnURL, "")
	}
}

func handleLoginPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := manager.GetInstance().SessionStore.Login(w, r)
		if err != nil {
			// always log the error
			logger.Errorf("Error logging in: %v from IP: %s", err, r.RemoteAddr)
		}

		var invalidCredentialsError *session.InvalidCredentialsError

		if errors.As(err, &invalidCredentialsError) {
			http.Error(w, "Username or password is invalid", http.StatusUnauthorized)
			return
		}

		if err != nil {
			// don't expose the error to the user
			http.Error(w, "An unexpected error occurred. See logs", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := manager.GetInstance().SessionStore.Logout(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// redirect to the login page if credentials are required
		prefix := getProxyPrefix(r)
		if config.GetInstance().HasCredentials() {
			http.Redirect(w, r, prefix+loginEndpoint, http.StatusFound)
		} else {
			http.Redirect(w, r, prefix+"/", http.StatusFound)
		}
	}
}
