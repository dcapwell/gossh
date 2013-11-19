package gossh

import (
	"fmt"
	"github.com/dcapwell/gossh/workpool"
	"log"
	"time"
)

const (
	MIN_POOL_SIZE = 1
	MAX_POOL_SIZE = 100
)

type Options struct {
	Concurrent int
	User       string
	Identity   string
	Options    map[string]string
}

type SshResponse struct {
	Code   int
	Stdout string
	Stderr string
}

type SshResponseContext struct {
	Hostname string
	Duration time.Duration
	Response SshResponse
}

type SshResponses struct {
	Responses []SshResponseContext
}

type Ssh interface {
	Run(hosts []string, cmd string, options Options) (SshResponses, error)
}

func NewSsh() Ssh {
	pool, _ := workpool.NewWorkPoolWithMax(MAX_POOL_SIZE)
	// error is only returned if max is not positive, so can ignore it
	return &sshProcessImpl{pool: pool}
}

type sshProcessImpl struct {
	pool workpool.WorkPool
}

//TODO should this return a chan?  WorkPool returns a chan, and can convert on the go.  encode/json doesn't support channels, so couldn't use this in http code then.
func (s *sshProcessImpl) Run(hosts []string, cmd string, options Options) (SshResponses, error) {
	// find how many hosts to run concurrently
	conc := runConcurrency(options, len(hosts))
	// create ssh worker per host, send to pool
	tasks := createTasks(hosts, cmd, options)
	chanResults, err := s.pool.Run(tasks, MIN_POOL_SIZE, conc)
	if err != nil {
		return SshResponses{}, err
	}
	rsp := waitForCompletion(chanResults, len(hosts))
	// wait for completion
	return rsp, nil
}

func waitForCompletion(results chan workpool.TaskResult, expectedResponses int) SshResponses {
	responses := make([]SshResponseContext, expectedResponses)
	idx := 0
	var err error
	for result := range results {
		responses[idx], err = taskResultToContext(result)
		if err != nil {
			// should this be ignored?
			log.Printf("[WARNING]\t%v", err)
		}
	}
	return SshResponses{responses}
}

func taskResultToContext(result workpool.TaskResult) (SshResponseContext, error) {
	if result.Status == workpool.SUCCESS {
		/*
		   switch rs := result.Result.(type) {
		   case SshResponseContext:
		     return rs, nil
		   }
		*/
		rs, ok := result.Result.(SshResponseContext)
		if ok {
			return rs, nil
		}
		return SshResponseContext{}, fmt.Errorf("Unable to convert result %v into SshResponseContext\n", result)
	}
	// ssh process task should not fail.  Error shouldn't be returned since ssh context contains errors as well.
	return SshResponseContext{}, fmt.Errorf("Unable to convert TaskResult %v to SshResponseContext; status is not supported %d\n", result, result.Status)
}

func createTasks(hosts []string, cmd string, options Options) chan workpool.Task {
	tasks := make(chan workpool.Task, len(hosts))
	go func() {
		for _, host := range hosts {
			task := NewSshTask(host, cmd, options)
			tasks <- task
		}
		close(tasks)
	}()
	return tasks
}

func NewSshTask(host string, cmd string, opt Options) workpool.Task {
	// use this method to switch impls
	return NewSshProcessTask(host, cmd, opt).Run
}

func runConcurrency(options Options, numHosts int) int {
	conc := MAX_POOL_SIZE
	if options.Concurrent > 0 {
		conc = options.Concurrent
	}
	if numHosts < conc {
		conc = numHosts
	}
	return conc
}
