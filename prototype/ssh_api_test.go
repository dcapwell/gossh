package prototype

import (
  "testing"
  "time"
)

type SshOptionsV1 struct {
  Concurrent  int
  User        string
  Identity    string
  Options     map[string]string
}

type SshResponseV1 struct {
  Code      int
  stdout    string
  stderr    string
}

type SshResponseContextV1 struct {
  Hostname        string
  Duration        time.Duration
  Response        SshResponseV1
}

type SshResponsesV1 struct {
  Responses     []SshResponseContextV1
}

type SshApiV1 interface {
  Run(hosts string, cmd string, options SshOptionsV1) (SshResponsesV1, error)
}

type SshApiV1Impl struct {}

func (s *SshApiV1Impl) Run(hosts string, cmd string, options SshOptionsV1) (SshResponsesV1, error) {
  return SshResponsesV1{}, nil
}

func TestApiV1(t *testing.T) {
  ssh := new(SshApiV1Impl)

  rsp, err := ssh.Run("example-[1-10]", "date", SshOptionsV1{
    Concurrent: 10,
    User: "techops",
  })
  if err != nil {
    t.Errorf("Error running ssh: %v\n", err)
    t.FailNow()
  }
  for _ = range rsp.Responses {
    // handle rsp
  }
}

// what if I want to run on multiple hosts, but no real pattern is able to limit to the hosts I want?
type SshApiV2 interface {
  Run(hosts []string, cmd string, options SshOptionsV1) (SshResponsesV1, error)
}

type SshApiV2Impl struct {
  maxConcurrency    int
}

func (s *SshApiV2Impl) Run(hosts []string, cmd string, options SshOptionsV1) (SshResponsesV1, error) {
  return SshResponsesV1{}, nil
}

func NewSshApiV2() SshApiV2 {
  /*
  ssh := new(SshApiV2Impl)
  ssh.maxConcurrency = 40
  return ssh
  */
  ssh := SshApiV2Impl{maxConcurrency:40}
  return &ssh
}

func TestApiV2(t *testing.T) {
  ssh := new(SshApiV2Impl)

  rsp, err := ssh.Run([]string{"example-1", "foo-1"}, "date", SshOptionsV1{
    Concurrent: 10,
    User: "techops",
  })
  if err != nil {
    t.Errorf("Error running ssh: %v\n", err)
    t.FailNow()
  }
  for _ = range rsp.Responses {
    // handle rsp
  }
}

func TestApiV2FromFactory(t *testing.T) {
  ssh := NewSshApiV2()
  t.Logf("Ssh %v\n", ssh)

  rsp, err := ssh.Run([]string{"example-1", "foo-1"}, "date", SshOptionsV1{
    Concurrent: 10,
    User: "techops",
  })
  if err != nil {
    t.Errorf("Error running ssh: %v\n", err)
    t.FailNow()
  }
  for _ = range rsp.Responses {
    // handle rsp
  }
}

func Explode(hosts string) []string {
  return []string{"host-1", "host-2"}
}

func TestApiV2ExplodeHosts(t *testing.T) {
  ssh := new(SshApiV2Impl)

  rsp, err := ssh.Run(Explode("example-[1..10]"), "date", SshOptionsV1{
    Concurrent: 10,
    User: "techops",
  })
  if err != nil {
    t.Errorf("Error running ssh: %v\n", err)
    t.FailNow()
  }
  for _ = range rsp.Responses {
    // handle rsp
  }
}
