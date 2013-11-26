package gossh

import (
  "fmt"
  "os"
  "io"
  "io/ioutil"
  "path/filepath"
  "code.google.com/p/go.crypto/ssh"
)

// create ssh task

// pulled from go.crytpo/ssh/test/test_unix_test.go
type keychain struct {
  keys []ssh.Signer
}

func newKeychain() (kc *keychain, err error) {
  kc = new(keychain)
  err = kc.load()
  return
}

func newKeychainWithKeys(keys ...string) (*keychain, error) {
  kc := new(keychain)
  for _, key := range keys {
    err := kc.loadPEM(key)
    if err != nil {
      return nil, err
    }
  }
  return kc, nil
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
