package main

import (
  "github.com/dcapwell/gossh"
  "github.com/dcapwell/gossh/utils"
  "flag"
  "io/ioutil"
  "os"
  "fmt"
  "log"
  "strings"
)

const (
  DEFAULT_CONCURRENT = 40
)

const USEAGE = `
gossh <hosts> [-n <num>] [-l <user>] [-i <path>] [-o Option=val] cmd
`

type Args struct {
  Hosts   string
  Cmd     string
  Opts    gossh.Options
}

func main() {
  args, err := parseArgs()
  if err != nil {
    panic(err)
  }
  ssh := gossh.NewSsh()
  hosts, err := utils.Expand(args.Hosts)
  if err != nil {
    panic(err)
  }
  rsp, err := ssh.Run(hosts, args.Cmd, args.Opts)
  if err != nil {
    panic(err)
  }
  for ctx := range rsp.Responses {
    log.Printf("Hostname: %s\nStdout: %s\nStderr: %s\n", ctx.Hostname, ctx.Response.Stdout, ctx.Response.Stderr)
  }
}

func parseArgs() (Args, error) {
  // $0, host, cmd
  if len(os.Args) < 3 {
    return Args{}, fmt.Errorf("%s\n", USEAGE)
  }

	flags := flag.NewFlagSet("Gossh", flag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)

  numConcurrent := flags.Int("n", DEFAULT_CONCURRENT, "defines how many concurrent resources to use")
  user := flags.String("l", "", "defines user to login as")
  identity := flags.String("i", "", "defines private key to use to login with")
  opts := new(OptArg)
  flags.Var(opts, "o", "ssh options")

  if err := flags.Parse(os.Args[2:]); err != nil {
    return Args{}, fmt.Errorf("%s\n", USEAGE)
  }

  args := Args{
    Hosts: os.Args[1],
    Cmd: strings.Join(flags.Args(), " "),
    Opts: gossh.Options{},
  }

  if numConcurrent != nil {
    args.Opts.Concurrent = *numConcurrent
  }
  if user != nil {
    args.Opts.User = *user
  }
  if identity != nil {
    args.Opts.Identity = *identity
  }
  optData, err := opts.ToMap()
  if err != nil {
    return Args{}, fmt.Errorf("%s\n", USEAGE)
  }
  args.Opts.Options = optData

  return args, nil
}

type OptArg struct {
  opts  []string
}

func (o *OptArg) String() string {
  return strings.Join(o.opts, ",")
}
func (o *OptArg) Set(opt string) error {
  o.opts = append(o.opts, opt)
  return nil
}

func (o *OptArg) ToMap() (map[string]string, error) {
  data := make(map[string]string)
  for _, opt := range o.opts {
    split := strings.Split(opt, "=")
    if len(split) != 2 {
      return nil, fmt.Errorf("Unable to parse option %v\n", opt)
    }
    data[split[0]] = split[1]
  }
  return data, nil
}
