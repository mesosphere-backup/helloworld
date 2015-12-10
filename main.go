package main

import (
  "flag"
  "fmt"
  "os"
  "io"
  "net/http"
)

var listener string

func init() {
  flag.StringVar(&listener, "listener", ":1234", "Listener address")
  flag.Parse()
}

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
  hostname, _ := os.Hostname()
  pid := os.Getpid()
  ret := fmt.Sprintf("Hello world (%d) app running at: %s, on hostname: %s\n", pid, listener, hostname)
  io.WriteString(w, ret)
}

func main() {
  fmt.Println("Hello world server starting up")
  http.HandleFunc("/", HelloServer)
  http.ListenAndServe(listener, nil)
}
