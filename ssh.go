package gossh

import (
  "time"
  "github.com/dcapwell/gossh/workpool"
)

type Options struct {
  Concurrent  int
  User        string
  Identity    string
  Options     map[string]string
}

type SshResponse struct {
  Code      int
  Stdout    string
  Stderr    string
}

type SshResponseContext struct {
  Hostname        string
  Duration        time.Duration
  Response        SshResponse
}

type SshResponses struct {
  Responses     []SshResponseContext
}

type Ssh interface {
  Run(hosts []string, cmd string, options Options) (SshResponses, error)
}

type sshImpl struct {
  pool      workpool.WorkPool
}

func (s *sshImpl) Run(hosts []string, cmd string, options Options) (SshResponses, error) {
  // find how many hosts to run concurrently
  // create new bounded worker pool
  // create ssh worker per host, send to pool
  // wait for completion
  return SshResponses{}, nil
}

func NewSsh() Ssh {
  pool, _ := workpool.NewWorkPoolWithMax(40)
  // error is only returned if max is not positive, so can ignore it
  return &sshImpl{pool: pool}
}


