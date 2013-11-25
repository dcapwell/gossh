package main

import (
  "code.google.com/p/gorest"
  "github.com/dcapwell/gossh"
  "github.com/dcapwell/gossh/utils"
  "log"
  "net/http"
  "fmt"
  "flag"
  "os"
  "io/ioutil"
  "strconv"
  "strings"
)

type SshService struct {
  gorest.RestService `root:"/ssh/"`

  hosts   gorest.EndPoint `method:"GET" path:"/{hosts:string}/" output:"SshCmdResponses"`
}

type SshCmdResponses struct {
  Responses []gossh.SshResponseContext
}

func (serv SshService) Hosts(hosts string) SshCmdResponses {
  // check for allowed query params
  cmd, opts, err := serv.parseQuery()
  if err != nil {
    return serv.badResponse(err)
  }
  ssh := gossh.NewSsh()
  h, err := utils.Expand(hosts)
  if err != nil {
    return serv.badResponse(err)
  }
  rsp, err := ssh.Run(h, cmd, opts)
  if err != nil {
    return serv.badResponse(err)
  }
  return SshCmdResponses{waitFor(rsp.Responses)}
}

func waitFor(ch chan gossh.SshResponseContext) []gossh.SshResponseContext {
  data := make([]gossh.SshResponseContext, 0)

  for d := range ch {
    data = append(data, d)
  }

  return data
}

func (serv SshService) parseQuery() (cmd string, opts gossh.Options, err error) {
  // map of string => []string
  val := serv.Context.Request().URL.Query()
  opts = gossh.Options{}
  // only first cmd is supported.  If user gives multiple, then it will be ignored
  // if not defined, will be empty
  cmd = val.Get("cmd")
  if cmd == "" {
    err = fmt.Errorf("cmd param not defined; cmd is a required param: %v", val)
    return
  }
  opts.User = val.Get("user")
  opts.Identity = val.Get("identity")
  // get concurrency
  concStr := val.Get("concurrent")
  if concStr != "" {
    conc, err := strconv.Atoi(concStr)
    if err != nil {
      //TODO if i don't add the outputs, compiler says that err is shadowed... find what that means
      return cmd, opts, err
    }
    opts.Concurrent = conc
  }
  // parse options
  o := val["option"]
  opts.Options = make(map[string]string)
  if o != nil {
    for _, option := range o {
      split := strings.Split(option, "=")
      if len(split) != 2 {
        err = fmt.Errorf("Option %s must contain exactly one '='; %v", option, val)
        return
      }
      opts.Options[split[0]] = split[1]
    }
  }
  return
}

func (serv SshService) badResponse(err error) SshCmdResponses {
  serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride([]byte(fmt.Sprintf("%v", err)))
  return SshCmdResponses{}
}

func main() {
  flags := flag.NewFlagSet("Gossh Agent", flag.ContinueOnError)
  flags.SetOutput(ioutil.Discard)
  base := "/"
  // looks like gorest doesn't work well unless at /
  // flags.StringVar(&base, "base", "/", "base path for http requests; defaults to '/'")
  port := ":7654"
  flags.StringVar(&port, "port", ":7654", "port to listen on; defaults to ':7654'")
  if err := flags.Parse(os.Args[1:]); err != nil {
    log.Fatal(err)
  }

  // register with gorest
  gorest.RegisterService(new(SshService))

  http.Handle(base,gorest.Handle())
  log.Fatal(http.ListenAndServe(port, nil))
}
