package yahooapi

import (
  "github.com/gorilla/mux"

)

func (a *YahooConfig) RegisterRoutes(r *mux.Router) {
  // auth routes
	r.HandleFunc("/yahoo/auth/", a.AuthYahoo)
	r.HandleFunc("/yahoo/auth/callback", a.AuthYahooCallback)

  // fantasy sports routes
  r.HandleFunc("/yahoo/users/", a.UserCollectionHandler)
}
