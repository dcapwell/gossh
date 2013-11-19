package workpool_test

import (
	. "github.com/dcapwell/gossh/workpool"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var NoOp = func() (interface{}, error) {
	return "NoOp", nil
}

func NoOpTasks(num int) chan Task {
	ch := make(chan Task, num)
	go func() {
		for i := 0; i < num; i++ {
			ch <- NoOp
		}
	}()
	return ch
}

var _ = Describe("Workpool", func() {

	Describe("create workpool", func() {
		Context("max resources", func() {
			It("is negative", func() {
				pool, err := NewWorkPoolWithMax(-5)
				Expect(err).To(HaveOccured())
				Expect(pool).To(BeNil())
			})
			It("is zero", func() {
				pool, err := NewWorkPoolWithMax(0)
				Expect(err).To(HaveOccured())
				Expect(pool).To(BeNil())
			})
			It("is positive", func() {
				pool, err := NewWorkPoolWithMax(5)
				Expect(err).To(BeNil())
				Expect(pool).ShouldNot(BeNil())
			})
		})
	})

	Describe("run no-op tasks", func() {
		Context("with", func() {
			It("1 resource", func() {
				pool, err := NewWorkPoolWithMax(1)
				Expect(err).To(BeNil())
				rsp, err := pool.Run(NoOpTasks(1), 1, 1)
				Expect(err).To(BeNil())
				Expect(rsp).ShouldNot(BeNil())
			})
			It("10 resource", func() {
			})
			It("100 resource", func() {
			})
		})
	})

})
