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
  for _, ctx := range rsp.Responses {
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

  return args, nil
}
