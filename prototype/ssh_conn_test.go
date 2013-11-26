package prototype

import (
  "testing"
  "fmt"
  "os"
  "io"
  "io/ioutil"
  "bytes"
  "path/filepath"
  "log"
  "code.google.com/p/go.crypto/ssh"
)

func TestSshCmd(t *testing.T) {
  kc := new(keychain)
  kc.load()
  config := &ssh.ClientConfig{
      User: os.Getenv("USER"),
      Auth: []ssh.ClientAuth{
          ssh.ClientAuthKeyring(kc),
      },
  }
  client, err := ssh.Dial("tcp", "localhost:22", config)
  if err != nil {
      panic("Failed to dial: " + err.Error())
  }

  // Each ClientConn can support multiple interactive sessions,
  // represented by a Session.
  session, err := client.NewSession()
  if err != nil {
      panic("Failed to create session: " + err.Error())
  }
  defer session.Close()

  // Once a Session is created, you can execute a single command on
  // the remote side using the Run method.
  var b bytes.Buffer
  session.Stdout = &b
  if err := session.Run("/usr/bin/whoami"); err != nil {
      panic("Failed to run: " + err.Error())
  }
  log.Printf("Result of running whoami via ssh: %s\n", b)
  //fmt.Println(b.String())
}

// pulled from go.crytpo/ssh/test/test_unix_test.go
type keychain struct {
  keys []ssh.Signer  
}

func (k *keychain) Key(i int) (ssh.PublicKey, error) {
  if i < 0 || i >= len(k.keys) {
    return nil, nil
  }
  return k.keys[i].PublicKey(), nil
}

func (k *keychain) Sign(i int, rand io.Reader, data []byte) (sig []byte, err error) {
  return k.keys[i].Sign(rand, data)
}

func (k *keychain) load() error {
  sshDir := fmt.Sprintf("%s/.ssh", os.Getenv("HOME"))
  return filepath.Walk(sshDir, func(path string, info os.FileInfo, err error) error {
    if err == nil {
      // no error reading file info, so lets see if its a pem
      // if not, just skip
      k.loadPEM(path)
    }
    return nil
  })
}

func (k *keychain) loadPEM(file string) error {
  buf, err := ioutil.ReadFile(file)
  if err != nil {
    return err  
  }
  key, err := ssh.ParsePrivateKey(buf)
  if err != nil {
    return err  
  }
  k.keys = append(k.keys, key)
  return nil
}
