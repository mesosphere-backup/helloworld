package main

import (
  "flag"
  "fmt"
  "os"
  "io"
  "net"
  "net/http"
  "time"
)

var healthy bool

var ln net.Listener

var listener string
var sayString string
func init() {
  flag.StringVar(&listener, "listener", ":1234", "Listener address")
  flag.StringVar(&sayString, "say-string", "Default Say String", "Say String")
  flag.Parse()
}

func HealthyServer(w http.ResponseWriter, req *http.Request) {
  if !healthy {
    w.WriteHeader(http.StatusInternalServerError)
  }
  newvalue := req.PostFormValue("healthy")
  if newvalue == "false" {
    healthy = false
  } else if newvalue == "true" {
    healthy = true
  }
  ret := fmt.Sprintf("Healthy: %v", healthy)
  io.WriteString(w, ret)
}

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
  hostname, _ := os.Hostname()
  pid := os.Getpid()
  ret := fmt.Sprintf("Hello world (%d) app running at: %s, on hostname: %s, and say string `%s` healthy: %v\n", pid, listener, hostname, sayString, healthy)
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

func DenyConnections(w http.ResponseWriter, req *http.Request) {
  ln.Close()
  ret := fmt.Sprintf("Denying future connections")
  io.WriteString(w, ret)
}

func main() {
  healthy = true
  var err error
  fmt.Println("Hello world server starting up")
  ln, err = net.Listen("tcp", listener)
  if err != nil {
    fmt.Println("Error, could not start listener", err)
    os.Exit(1)
  }
  http.HandleFunc("/", HelloServer)
  http.HandleFunc("/stream", StreamingServer)
  http.HandleFunc("/deny", DenyConnections)
  http.HandleFunc("/healthy", HealthyServer)
  http.Serve(ln, nil)
  select{}
}
