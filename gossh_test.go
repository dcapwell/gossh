package gossh_test

import (
	. "github.com/dcapwell/gossh"
	"github.com/dcapwell/gossh/workpool"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
)

var _ = Describe("Gossh", func() {
	ssh := NewSsh()

	Describe("run multiple ssh commands concurrently", func() {
		Context("with one host", func() {
			Context("and valid command", func() {
				rsp, err := ssh.Run([]string{"localhost"}, "date", Options{})

				It("should succeed", func() {
					Expect(err).To(BeNil())
				})

				It("with one result", func() {
					Expect(rsp.Responses).To(HaveLen(1))
					Expect(rsp.Responses[0].Hostname).To(Equal("localhost"))
					Expect(rsp.Responses[0].Response.Code).To(Equal(workpool.SUCCESS))
					//TODO make this work on any envo
					Expect(rsp.Responses[0].Response.Stdout).To(ContainSubstring("PST"))
				})
			})
			Context("with invalid command", func() {
				rsp, err := ssh.Run([]string{"localhost"}, "thiscmdreallyshouldntexist", Options{})

				It("should not return error", func() {
					Expect(err).To(BeNil())
				})
				It("should have only one response", func() {
					Expect(rsp.Responses).To(HaveLen(1))
				})
				It("should have response from host", func() {
					localRsp := rsp.Responses[0]
					Expect(localRsp.Hostname).To(Equal("localhost"))
				})
				It("should have failed", func() {
					localRsp := rsp.Responses[0]
					Expect(localRsp.Response).ShouldNot(BeNil())
					Expect(localRsp.Response.Code).To(BeGreaterThan(0))
				})
			})
		})
		Context("with mulitple hosts", func() {
			Context("with valid command", func() {
				rsp, err := ssh.Run([]string{"localhost", "localhost"}, "date", Options{})
				It("should not return error", func() {
					Expect(err).To(BeNil())
				})
				It("should have only two response", func() {
					Expect(rsp.Responses).To(HaveLen(2))
				})
				It("should have response from host", func() {
					localRsp := rsp.Responses[0]
					Expect(localRsp.Hostname).To(Equal("localhost"))
					localRsp = rsp.Responses[1]
					Expect(localRsp.Hostname).To(Equal("localhost"))
				})
				It("should have successed", func() {
					Expect(rsp.Responses[0].Response.Code).To(Equal(workpool.SUCCESS))
					Expect(rsp.Responses[0].Response.Stdout).To(ContainSubstring("PST"))

					Expect(rsp.Responses[1].Response.Code).To(Equal(workpool.SUCCESS))
					Expect(rsp.Responses[1].Response.Stdout).To(ContainSubstring("PST"))
				})
			})
			Context("with invalid command", func() {
				rsp, err := ssh.Run([]string{"localhost", "localhost"}, "thiscmdreallyshouldntexist", Options{})
				log.Printf("Response from multi requests, single host, invalid cmd: %v\n", rsp)

				It("should not return error", func() {
					Expect(err).To(BeNil())
				})
				It("should have only two response", func() {
					Expect(rsp.Responses).To(HaveLen(2))
				})
				It("should have response from host", func() {
					Expect(rsp.Responses[0].Hostname).To(Equal("localhost"))
					Expect(rsp.Responses[1].Hostname).To(Equal("localhost"))
				})
				It("should have failed", func() {
					localRsp := rsp.Responses[0]
					Expect(localRsp.Response).ShouldNot(BeNil())
					Expect(localRsp.Response.Code).To(BeGreaterThan(0))

					localRsp = rsp.Responses[1]
					Expect(localRsp.Response).ShouldNot(BeNil())
					Expect(localRsp.Response.Code).To(BeGreaterThan(0))
				})
			})
		})
	})
})
