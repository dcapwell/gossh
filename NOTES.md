Functional interface for ssh

Run command on host
```
import "github.com/dcapwell/gossh"
rsp, err := gossh.Run("localhost", "date", gossh.Config{
  User: "randomuser",
  Identity: "/home/randomuser/.ssh/id_dsa",
  Native: false,
  Options: gossh.Opts().
    Add("PreferredAuthentications", "publickey").
    Add("StrictHostKeyChecking", "no")
  ),
})
if err != nil {
  log.Fatalf("Unable to run ssh command: %v\n", err)
}
log.Printf("%s\n", json.Marshal(rsp)) // wont compile, but you get the idea
/*
output:
{
  "Code": 0,
  "Stdout": "Wed Nov 27 09:37:29 PST 2013\n",
  "Stderr": ""
}
*/
```

Run command on hosts concurrently
```
import "github.com/dcapwell/gossh"
import "github.com/dcapwell/gossh/util"
rsp, err := gossh.RunMulti(util.Expand("localhost", "localhost"), "date", gossh.Config{
  User: "randomuser",
  Identity: "/home/randomuser/.ssh/id_dsa",
  Native: false,
  Options: gossh.Opts().
    Add("PreferredAuthentications", "publickey").
    Add("StrictHostKeyChecking", "no")
  ),
})
if err != nil {
  log.Fatalf("Unable to run ssh command: %v\n", err)
}
log.Printf("%s\n", json.Marshal(rsp)) // wont compile, but you get the idea
/*
output:
{
  "Code": 0,
  "Stdout": "Wed Nov 27 09:37:29 PST 2013\n",
  "Stderr": ""
}
*/
```

Run command on host, and return debug information
```
import "github.com/dcapwell/gossh"
rsp, err := gossh.RunWithDebug("localhost", "date", gossh.Config{
  User: "randomuser",
  Identity: "/home/randomuser/.ssh/id_dsa",
  Native: false,
  Options: gossh.Opts().
    Add("PreferredAuthentications", "publickey").
    Add("StrictHostKeyChecking", "no")
  ),
})
if err != nil {
  log.Fatalf("Unable to run ssh command: %v\n", err)
}
log.Printf("%s\n", json.Marshal(rsp)) // wont compile, but you get the idea
/*
output:
{
  "Debug": {
    "Hostname": "localhost",
    "Duration": "82ms",
    "Config": {
      "User": "randomuser",
      "Identity": "/home/randomuser/.ssh/id_dsa",
      "Native": false,
      "Options": {
        "PreferredAuthentications": "publickey",
        "StrictHostKeyChecking": "no"
      }
    }
  },
  "Response" : {
    "Code": 0,
    "Stdout": "Wed Nov 27 09:37:29 PST 2013\n",
    "Stderr": ""
  }
}
*/
```
