package prototype

import (
  "code.google.com/p/gorest"
  "log"
  "net/http"
)

type SshService struct {
  gorest.RestService `root:"/ssh/"`

  helloWorld  gorest.EndPoint `method:"GET" path:"/hello-world/" produces:"shouldfreakout" output:"string"`
  helloWorldTwo  gorest.EndPoint `method:"GET" path:"/hello-world-2/" output:"string"`
  serde       gorest.EndPoint `method:"GET" path:"/serde/" output:"map[string]string"`
}

func (serv SshService) HelloWorld() string {
  return "Hello World!"
}

func (serv SshService) HelloWorldTwo() string {
  serv.ResponseBuilder().
    AddHeader("Access-Control-Allow-Origin","http://127.0.0.1:8888").
    AddHeader("Access-Control-Allow-Headers","X-HTTP-Method-Override").
    AddHeader("Access-Control-Allow-Headers","X-Xsrf-Cookie").
    AddHeader("Access-Control-Expose-Headers","X-Xsrf-Cookie").
    SetResponseCode(200).
    SetContentType(gorest.Text_RichText)
  return "Hello World! Two!"
}

func (serv SshService) Serde() map[string]string {
  data := make(map[string]string)
  data["Hello"] = "World!"
  return data
}

func (serv SshService) Builder() {
  rb := serv.RB()
  rb.Write([]byte("Hello World"))
}

func main() {
  // register with gorest
  gorest.RegisterService(new(SshService))

  http.Handle("/",gorest.Handle())
  log.Fatal(http.ListenAndServe(":7654", nil))
}
