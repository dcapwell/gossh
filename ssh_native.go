package gossh

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// create a new native ssh task
func newSshNativeTask(host string, cmd string, opt Options) func() (interface{}, error) {
	state := &sshNativeTask{
		Host: host,
		Cmd:  cmd,
		Opts: opt,
	}
	return state.run
}

type sshNativeTask struct {
	Host string
	Cmd  string
	Opts Options
}

// workpool task function.  Runs ssh cmd
func (s *sshNativeTask) run() (interface{}, error) {
	client, err := s.newClient()
	// create session
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// run cmd
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	ctx := createContext(s.Host)
	if err := session.Run(s.Cmd); err != nil {
		// if of type ExitError, then its a remote issue, add to code
		waitmsg, ok := err.(*ssh.ExitError)
		if ok {
			ctx.Response.Code = waitmsg.ExitStatus()
		} else {
			// else return err
			return ctx, err
		}
	}

	ctx.Response.Stdout = stdout.String()
	ctx.Response.Stderr = stderr.String()

	return ctx, nil
}

func (s *sshNativeTask) newClient() (*ssh.ClientConn, error) {
	config, err := s.newClientConfig()
	if err != nil {
		return nil, err
	}
	// create client
	//TODO allow user to override port
	client, err := ssh.Dial("tcp", s.Host+":22", config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *sshNativeTask) newClientConfig() (*ssh.ClientConfig, error) {
	// get auth
	// create keychain; only keys are supported right now
	var kc *keychain
	var err error
	if s.Opts.Identity != "" {
		kc, err = newKeychainWithKeys(s.Opts.Identity)
	} else {
		kc, err = newKeychain()
	}
	if err != nil {
		return nil, err
	}
	// get user
	user := s.Opts.User
	if user == "" {
		user = os.Getenv("USER")
	}
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.ClientAuth{
			ssh.ClientAuthKeyring(kc),
		},
	}
	return config, nil
}

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
