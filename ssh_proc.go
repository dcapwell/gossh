package gossh

import (
  "os/exec"
  "syscall"
  "io/ioutil"
  "fmt"
)

// this file is for SshTask that calls the local ssh shell command.
type SshProcessTask struct {
	Host    string
	Cmd     string
	Options Options
}

func (s *SshProcessTask) Run() (interface{}, error) {
	// must return of type (SshResponseContext, error)
  cmd := exec.Command("/usr/bin/ssh", s.Host, s.Cmd)
  stdout, err := cmd.StdoutPipe()
  if err != nil {
    return nil, err
  }
  stderr, err := cmd.StderrPipe()
  if err != nil {
    return nil, err
  }

  if err := cmd.Start(); err != nil {
    return nil, err
  }

  ctx := createContext(s.Host)

  // if this is after wait, it seems that the stream is closed
  // so no longer readable
  out, _ := ioutil.ReadAll(stdout)
  ctx.Response.Stdout = string(out)
  outerr, _ := ioutil.ReadAll(stderr)
  ctx.Response.Stderr = string(outerr)

  err = cmd.Wait()

  exitCode, err := exitCode(err)
  if err != nil {
    // run on non supported OS
    return ctx, err
  }
  ctx.Response.Code = exitCode

	return ctx, nil
}

func createContext(host string) SshResponseContext {
  rsp := SshResponse{}
  ctx := SshResponseContext{
    Hostname: host,
    Response: rsp,
  }
  return ctx
}

func exitCode(err error) (int, error) {
  if err != nil {
    // it puts exit code in err... grrr
    exitErr, ok := err.(*exec.ExitError)
    if ok {
      sys := exitErr.Sys()
      // this is system dependent.  This is the unix way
      waitStatus, ok := sys.(syscall.WaitStatus)
      if ok {
        return waitStatus.ExitStatus(), nil
      }
    }
    return -1, fmt.Errorf("Unsupported OS; expected syscall status to be on a unix environment", err)
  }
  // was successful, so exit code is 0
  return 0, nil
}

func NewSshProcessTask(host string, cmd string, opt Options) *SshProcessTask {
	return &SshProcessTask{
		Host:    host,
		Cmd:     cmd,
		Options: opt,
	}
}
