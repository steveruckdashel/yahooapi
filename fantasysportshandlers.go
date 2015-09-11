package yahooapi

import (
  "net/http"
  "io"
  // "fmt"
  "encoding/json"
  "log"
)

func (y *YahooConfig) UserCollectionGamesHandler(w http.ResponseWriter, r *http.Request) {
  user := y.GetUserCollectionGames(r)
  //io.WriteString(w, user.Body)
  // io.WriteString(w, fmt.Sprintf("%v", user))
  b, err := json.MarshalIndent(user, "", "  ")
  if err != nil {
    log.Fatal(err)
  }
  w.Write(b)
}

func (y *YahooConfig) UserCollectionLeaguesHandler(w http.ResponseWriter, r *http.Request) {
  user := y.GetUserCollectionLeagues(r)
  io.WriteString(w, user.Body)
  // io.WriteString(w, fmt.Sprintf("%v", user))
  // b, err := json.MarshalIndent(user, "", "  ")
  // if err != nil {
  //   log.Fatal(err)
  // }
  // w.Write(b)
}

func (y *YahooConfig) UserCollectionTeamsHandler(w http.ResponseWriter, r *http.Request) {
  user := y.GetUserCollectionTeams(r)
  io.WriteString(w, user.Body)
  // io.WriteString(w, fmt.Sprintf("%v", user))
  // b, err := json.MarshalIndent(user, "", "  ")
  // if err != nil {
  //   log.Fatal(err)
  // }
  // w.Write(b)
}
