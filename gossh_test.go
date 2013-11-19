package gossh_test

import (
	. "github.com/dcapwell/gossh"
	"github.com/dcapwell/gossh/workpool"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gossh", func() {
	ssh := NewSsh()

	Describe("run multiple ssh commands concurrently", func() {
		Context("with one host", func() {
			rsp, err := ssh.Run([]string{"localhost"}, "date", Options{})

			It("should succeed", func() {
				Expect(err).To(BeNil())
			})

			It("with one result", func() {
				Expect(rsp.Responses).To(HaveLen(1))
				Expect(rsp.Responses[0].Hostname).To(Equal("localhost"))
				Expect(rsp.Responses[0].Response.Code).To(Equal(workpool.SUCCESS))
				Expect(rsp.Responses[0].Response.Stdout).To(ContainSubstring("PST"))
			})
		})
	})
})
