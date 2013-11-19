package gossh

// this file is for SshTask that calls the local ssh shell command.
type SshProcessTask struct {
	Host    string
	Cmd     string
	Options Options
}

func (s *SshProcessTask) Run() (interface{}, error) {
	// must return of type (SshResponseContext, error)
	return nil, nil
}

func NewSshProcessTask(host string, cmd string, opt Options) *SshProcessTask {
	return &SshProcessTask{
		Host:    host,
		Cmd:     cmd,
		Options: opt,
	}
}
