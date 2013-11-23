package main

import (
  "code.google.com/p/gorest"
  "github.com/dcapwell/gossh"
  "github.com/dcapwell/gossh/utils"
  "log"
  "net/http"
  "fmt"
)

type SshService struct {
  gorest.RestService `root:"/ssh/"`

  hosts   gorest.EndPoint `method:"GET" path:"/{hosts:string}/" output:"SshCmdResponses"`
}

type SshCmdResponses struct {
  Responses []gossh.SshResponseContext
}

func WaitFor(ch chan gossh.SshResponseContext) []gossh.SshResponseContext {
  data := make([]gossh.SshResponseContext, 0)

  for d := range ch {
    data = append(data, d)
  }

  return data
}

func (serv SshService) Hosts(hosts string) SshCmdResponses {
  ssh := gossh.NewSsh()
  h, err := utils.Expand(hosts)
  if err != nil {
    badResponse(serv, err)
    return SshCmdResponses{}
  }
  rsp, err := ssh.Run(h, "date", gossh.Options{})
  if err != nil {
    badResponse(serv, err)
    return SshCmdResponses{}
  }
  /*
  for ctx := range rsp.Responses {
    log.Printf("Hostname: %s\nStdout: %s\nStderr: %s\n", ctx.Hostname, ctx.Response.Stdout, ctx.Response.Stderr)
  }
  */
  //return rsp
  return SshCmdResponses{WaitFor(rsp.Responses)}
}

func badResponse(serv SshService, err error) {
  serv.ResponseBuilder().SetResponseCode(400).WriteAndOveride([]byte(fmt.Sprintf("%v", err)))
}


func main() {
  // register with gorest
  gorest.RegisterService(new(SshService))

  http.Handle("/",gorest.Handle())
  log.Fatal(http.ListenAndServe(":7654", nil))
}
