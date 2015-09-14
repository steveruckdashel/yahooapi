package yahooapi

import (
	// "io"
	"net/http"
	// "fmt"
	"encoding/json"
	"log"
)

func (y *YahooConfig) UserCollectionGamesHandler(w http.ResponseWriter, r *http.Request) {
	user := y.GetUserCollectionGames(r)
	// io.WriteString(w, user.Body)
	// io.WriteString(w, fmt.Sprintf("%v", user))
	b, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(b)
}

func (y *YahooConfig) UserCollectionLeaguesHandler(w http.ResponseWriter, r *http.Request) {
	user := y.GetUserCollectionLeagues(r)
	// io.WriteString(w, user.Body)
	// io.WriteString(w, fmt.Sprintf("%v", user))
	b, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
	  log.Fatal(err)
	}
	w.Write(b)
}

func (y *YahooConfig) UserCollectionTeamsHandler(w http.ResponseWriter, r *http.Request) {
	user := y.GetUserCollectionTeams(r)
	// io.WriteString(w, user.Body)
	// io.WriteString(w, fmt.Sprintf("%v", user))
	b, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
	  log.Fatal(err)
	}
	w.Write(b)
}

func (y *YahooConfig) UserCollectionAllHandler(w http.ResponseWriter, r *http.Request) {
	user := y.GetUserCollectionAll(r)
	// io.WriteString(w, user.Body)
	// io.WriteString(w, fmt.Sprintf("%v", user))
	b, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
	  log.Fatal(err)
	}
	w.Write(b)
}

func (y *YahooConfig) LeagueScoreboardHandler(w http.ResponseWriter, r *http.Request) {
	scoreboard := y.GetLeagueScoreboard(r)
	// io.WriteString(w, user.Body)
	// io.WriteString(w, fmt.Sprintf("%v", user))
	b, err := json.MarshalIndent(scoreboard, "", "  ")
	if err != nil {
	  log.Fatal(err)
	}
	w.Write(b)
}

func (y *YahooConfig) LeagueStandingsHandler(w http.ResponseWriter, r *http.Request) {
	standings := y.GetLeagueStandings(r)
	// io.WriteString(w, user.Body)
	// io.WriteString(w, fmt.Sprintf("%v", user))
	b, err := json.MarshalIndent(standings, "", "  ")
	if err != nil {
	  log.Fatal(err)
	}
	w.Write(b)
}
