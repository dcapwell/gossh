package gossh_test

import (
	"github.com/dcapwell/gossh"
	"log"
)

func Example() {
	ssh := gossh.NewSsh()
	rsp, err := ssh.Run([]string{"localhost"}, "date", gossh.Options{})
	if err != nil {
		log.Fatalf("Unable to run command 'date' on host 'localhost': %v\n", err)
	}

	for _, ctx := range rsp.Responses {
		log.Printf("Response from host %s: %s\n", ctx.Hostname, ctx.Response.Stdout)
	}
}
