package yahooapi

import (
	"github.com/gorilla/mux"
)

func (a *YahooConfig) RegisterRoutes(r *mux.Router) {
	// auth routes
	r.HandleFunc("/yahoo/auth/", a.AuthYahoo)
	r.HandleFunc("/yahoo/auth/callback", a.AuthYahooCallback)

	// fantasy sports routes
	r.HandleFunc("/yahoo/users/games", a.UserCollectionGamesHandler)
	r.HandleFunc("/yahoo/users/game/{game_keys:[0-9]+}", a.UserCollectionAllHandler)
	r.HandleFunc("/yahoo/users/game/{game_keys:[0-9]+}/leagues", a.UserCollectionLeaguesHandler)
	r.HandleFunc("/yahoo/users/game/{game_keys:[0-9]+}/teams", a.UserCollectionTeamsHandler)
  r.HandleFunc("/yahoo/users/leagues/{league_keys:[0-9a-zA-Z\\.]+}/scoreboard", a.LeagueScoreboardHandler)
  r.HandleFunc("/yahoo/users/leagues/{league_keys:[0-9a-zA-Z\\.]+}/standings", a.LeagueStandingsHandler)
}
