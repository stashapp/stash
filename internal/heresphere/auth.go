package heresphere

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/session"
)

/*
 * This auxiliary function finds if a login is needed, and auth is correct.
 */
func basicLogin(username string, password string) bool {
	// If needs creds, try login
	if config.GetInstance().HasCredentials() {
		err := manager.GetInstance().SessionStore.LoginPlain(username, password)
		return err != nil
	}
	return false
}

/*
 * This auxiliary function finds if the request has a valid auth token.
 */
func HeresphereHasValidToken(r *http.Request) bool {
	// Check header auth
	apiKey := r.Header.Get(HeresphereAuthHeader)

	// Check url query auth
	if len(apiKey) == 0 {
		apiKey = r.URL.Query().Get(session.ApiKeyParameter)
	}

	return len(apiKey) > 0 && apiKey == config.GetInstance().GetAPIKey()
}

/*
 * This auxiliary function adds an auth token to a url
 */
func addApiKey(urlS string) string {
	// Parse URL
	u, err := url.Parse(urlS)
	if err != nil {
		// shouldn't happen
		panic(err)
	}

	// Add apikey if applicable
	if config.GetInstance().GetAPIKey() != "" {
		v := u.Query()
		if !v.Has("apikey") {
			v.Set("apikey", config.GetInstance().GetAPIKey())
		}
		u.RawQuery = v.Encode()
	}

	return u.String()
}

/*
 * This auxiliary writes a library with a fake name upon auth failure
 */
func writeNotAuthorized(w http.ResponseWriter, r *http.Request, msg string) {
	// Banner
	banner := HeresphereBanner{
		Image: fmt.Sprintf("%s%s",
			manager.GetBaseURL(r),
			"/apple-touch-icon.png",
		),
		Link: fmt.Sprintf("%s%s",
			manager.GetBaseURL(r),
			"/",
		),
	}
	// Default video
	library := HeresphereIndexEntry{
		Name: msg,
		List: []string{},
	}
	// Index
	idx := HeresphereIndex{
		Access:  HeresphereBadLogin,
		Banner:  banner,
		Library: []HeresphereIndexEntry{library},
	}

	// Create a JSON encoder for the response writer
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(idx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
