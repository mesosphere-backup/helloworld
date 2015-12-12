package main

import (
  "flag"
  "fmt"
  "os"
  "io"
  "net/http"
  "time"
)

var listener string
var sayString string
func init() {
  flag.StringVar(&listener, "listener", ":1234", "Listener address")
  flag.StringVar(&sayString, "say-string", "Default Say String", "Say String")
  flag.Parse()
}

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
  hostname, _ := os.Hostname()
  pid := os.Getpid()
  ret := fmt.Sprintf("Hello world (%d) app running at: %s, on hostname: %s, and say string `%s`\n", pid, listener, hostname, sayString)
  io.WriteString(w, ret)
}

func StreamingServer(w http.ResponseWriter, req *http.Request) {
  f := w.(http.Flusher)
  notify := w.(http.CloseNotifier).CloseNotify()  
  hostname, _ := os.Hostname()
  pid := os.Getpid()
  count := 0
  outer:
  for {
    select {
      case <-notify: { break outer }
      case <-time.After(1 * time.Second): {
        ret := fmt.Sprintf("%d: Hello world (%d) app running at: %s, on hostname: %s, and say string `%s`\n", count, pid, listener, hostname, sayString)
        io.WriteString(w, ret)
        f.Flush()
      }
    }
    count = count + 1
  }
}

func main() {
  fmt.Println("Hello world server starting up")
  http.HandleFunc("/", HelloServer)
  http.HandleFunc("/stream", StreamingServer)
  http.ListenAndServe(listener, nil)
}
