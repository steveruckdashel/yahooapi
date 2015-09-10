package yahooapi

import (
  "net/http"
  "io"
)

func (y *YahooConfig) UserCollectionHandler(w http.ResponseWriter, r *http.Request) {
  user := y.GetUserCollection(r)
  io.WriteString(w, user.Body)
}
