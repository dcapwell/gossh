package gossh

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"syscall"
	"time"
)

// this file is for SshTask that calls the local ssh shell command.
type sshProcessTask struct {
	Host    string
	Cmd     string
	Options Options
}

func (s *sshProcessTask) generateCmdArguments() []string {
	// make a slice of init size 4, but can expand to 100
	cmd := make([]string, 0)

	cmd = append(cmd, "-n")

	// add user/identity
	if s.Options.User != "" {
		cmd = append(cmd, "-l", s.Options.User)
	}
	if s.Options.Identity != "" {
		cmd = append(cmd, "-i", s.Options.Identity)
	}

	// add options
	for key, value := range s.Options.Options {
		cmd = append(cmd, "-o", key+"="+value)
	}

	// add host
	cmd = append(cmd, s.Host)

	// last action, add cmd
	cmd = append(cmd, s.Cmd)

	return cmd
}

func (s *sshProcessTask) Run() (interface{}, error) {
	start := time.Now()
	// must return of type (SshResponseContext, error)
	//cmd := exec.Command("/usr/bin/ssh", s.Host, s.Cmd)
	cmd := exec.Command("ssh", s.generateCmdArguments()...)
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
	total := time.Now().Sub(start)
	ctx.Duration = fmt.Sprintf("%dms", total/time.Millisecond)
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

func newSshProcessTask(host string, cmd string, opt Options) *sshProcessTask {
	return &sshProcessTask{
		Host:    host,
		Cmd:     cmd,
		Options: opt,
	}
}
