package unix_socket

import (
  "net"
  "net/http"
  "testing"
  "time"
  "os"
)

const FILENAME = "/tmp/go-unix.sock"

func TestSocketExists(t *testing.T) {
  fd, err := net.Listen("unix", FILENAME)
  if err != nil {
    t.Fatal(err)
  }
  defer fd.Close()

  if _, err := os.Stat(FILENAME); os.IsNotExist(err) {
    t.Fatal(err)
  }
}

func handler(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/plain")
  w.Write([]byte("This is an example server.\n"))
}

/*
func close(l *net.Listener) func(http.ResponseWriter, *http.Request) {
  return func (w http.ResponseWriter, req *http.Request) {
    l.Close()
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte("closed listener.\n"))
  }
}
*/

func TestRunServer(t *testing.T) {
  fd, err := net.Listen("unix", FILENAME)
  if err != nil {
    t.Fatal(err)
  }
  go func() {
    time.Sleep(1 * time.Second)
    fd.Close()
  }()

  h := http.NewServeMux()
  h.HandleFunc("/", handler)

  err = http.Serve(fd, h)
  if err != nil {
    t.Log(err)
    if err.Error() != "accept unix "+FILENAME+": use of closed network connection" {
      t.Fatal(err)
    }
  }
}
