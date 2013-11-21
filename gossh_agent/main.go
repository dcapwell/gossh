package main

import (
  "log"
  "net/http"
  "fmt"
  "html"
)

func main() {
  http.HandleFunc("/ssh", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
  })

  log.Fatal(http.ListenAndServe(":7654", nil))
}
