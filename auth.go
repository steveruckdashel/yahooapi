package yahooapi

import (
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"net/url"
	"encoding/gob"
)

type YahooConfig struct {
	conf         *oauth2.Config
	SessionStore sessions.Store
	landing      string
}

func NewYahooConfig(clientID, clientSecret string, scopes []string, hostName string, landing string, sessionStore sessions.Store) *YahooConfig {
	gob.Register(&oauth2.Token{})

	return &YahooConfig{
		conf: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       scopes,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://api.login.yahoo.com/oauth2/request_auth",
				TokenURL: "https://api.login.yahoo.com/oauth2/get_token",
			},
			RedirectURL: hostName + "/yahoo/auth/callback",
		},
		SessionStore: sessionStore,
		landing:      landing,
	}
}

func (a *YahooConfig) AuthYahoo(w http.ResponseWriter, r *http.Request) {
	session, err := a.SessionStore.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	urlStr := a.conf.AuthCodeURL(session.Values["state"].(string), oauth2.AccessTypeOnline)
	urlStrUnesc, err := url.QueryUnescape(urlStr)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Visit the URL for the auth dialog: %v", urlStrUnesc)

	http.Redirect(w, r, urlStrUnesc, 302)
}

func (a *YahooConfig) AuthYahooCallback(w http.ResponseWriter, r *http.Request) {
	session, err := a.SessionStore.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// Use the authorization code that is pushed to the redirect URL.
	// NewTransportWithCode will do the handshake to retrieve
	// an access token and initiate a Transport that is
	// authorized and authenticated by the retrieved token.
	code := r.FormValue("code")

	tok, err := a.conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}
	session.Values["token"] = tok
	session.Values["xoauth_yahoo_guid"] = r.FormValue("xoauth_yahoo_guid")
	session.Save(r, w)

	// a.conf.Client(oauth2.NoContext, tok)

	http.Redirect(w, r, a.landing, 302)
}
