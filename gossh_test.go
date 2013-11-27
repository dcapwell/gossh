package gossh_test

import (
	. "github.com/dcapwell/gossh"
	"github.com/dcapwell/gossh/workpool"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
)

var _ = Describe("Gossh", func() {
  Describe("multi hosts", func() {
    ssh := NewSsh()

    Describe("run multiple ssh commands concurrently", func() {
      Context("with one host", func() {
        Context("and valid command", func() {
          rsp, err := ssh.Run([]string{"localhost"}, "date", Options{})
          data := chanToSlize(rsp.Responses)

          It("should succeed", func() {
            Expect(err).To(BeNil())
          })

          It("with one result", func() {
            Expect(data).To(HaveLen(1))
            Expect(data[0].Hostname).To(Equal("localhost"))
            Expect(data[0].Response.Code).To(Equal(workpool.SUCCESS))
            Expect(data[0].Response.Stdout).To(ContainSubstring("PST"))
          })
        })
        Context("with invalid command", func() {
          rsp, err := ssh.Run([]string{"localhost"}, "thiscmdreallyshouldntexist", Options{})
          data := chanToSlize(rsp.Responses)

          It("should not return error", func() {
            Expect(err).To(BeNil())
          })
          It("should have only one response", func() {
            Expect(data).To(HaveLen(1))
          })
          It("should have response from host", func() {
            localRsp := data[0]
            Expect(localRsp.Hostname).To(Equal("localhost"))
          })
          It("should have failed", func() {
            localRsp := data[0]
            Expect(localRsp.Response).ShouldNot(BeNil())
            Expect(localRsp.Response.Code).To(BeGreaterThan(0))
          })
        })
      })
      Context("with mulitple hosts", func() {
        Context("with valid command", func() {
          rsp, err := ssh.Run([]string{"localhost", "localhost"}, "date", Options{})
          data := chanToSlize(rsp.Responses)

          It("should not return error", func() {
            Expect(err).To(BeNil())
          })
          It("should have only two response", func() {
            Expect(data).To(HaveLen(2))
          })
          It("should have response from host", func() {
            localRsp := data[0]
            Expect(localRsp.Hostname).To(Equal("localhost"))
            localRsp = data[1]
            Expect(localRsp.Hostname).To(Equal("localhost"))
          })
          It("should have successed", func() {
            Expect(data[0].Response.Code).To(Equal(workpool.SUCCESS))
            Expect(data[0].Response.Stdout).To(ContainSubstring("PST"))

            Expect(data[1].Response.Code).To(Equal(workpool.SUCCESS))
            Expect(data[1].Response.Stdout).To(ContainSubstring("PST"))
          })
        })
        Context("with invalid command", func() {
          rsp, err := ssh.Run([]string{"localhost", "localhost"}, "thiscmdreallyshouldntexist", Options{})
          data := chanToSlize(rsp.Responses)

          log.Printf("Response from multi requests, single host, invalid cmd: %v\n", rsp)

          It("should not return error", func() {
            Expect(err).To(BeNil())
          })
          It("should have only two response", func() {
            Expect(data).To(HaveLen(2))
          })
          It("should have response from host", func() {
            Expect(data[0].Hostname).To(Equal("localhost"))
            Expect(data[1].Hostname).To(Equal("localhost"))
          })
          It("should have failed", func() {
            localRsp := data[0]
            Expect(localRsp.Response).ShouldNot(BeNil())
            Expect(localRsp.Response.Code).To(BeGreaterThan(0))

            localRsp = data[1]
            Expect(localRsp.Response).ShouldNot(BeNil())
            Expect(localRsp.Response.Code).To(BeGreaterThan(0))
          })
        })
      })
    })
  })
})

func chanToSlize(ch chan SshResponseContext) []SshResponseContext {
	data := make([]SshResponseContext, 0)
	for cxt := range ch {
		data = append(data, cxt)
	}
	return data
}
