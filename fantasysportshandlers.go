package yahooapi

import (
  "net/http"
  // "io"
  // "fmt"
  "encoding/json"
  "log"
)

func (y *YahooConfig) UserCollectionHandler(w http.ResponseWriter, r *http.Request) {
  user := y.GetUserCollection(r)
  //io.WriteString(w, user.Body)
  // io.WriteString(w, fmt.Sprintf("%v", user))
  b, err := json.MarshalIndent(user, "", "  ")
  if err != nil {
    log.Fatal(err)
  }
  w.Write(b)
}
