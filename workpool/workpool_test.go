package workpool_test

import (
	. "github.com/dcapwell/gossh/workpool"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
)

var NoOp = func() (interface{}, error) {
	return "NoOp", nil
}

type StatefulNoOp struct {}
func (s StatefulNoOp) Apply() (interface{}, error) {
  return "Stateful NoOP", nil
}
//var State = StatefulNoOp{}

func NoOpTasks(num int) chan Task {
	ch := make(chan Task, num)
	go func() {
		for i := 0; i < num; i++ {
      ch <- NoOp
      // ch <- State.Apply
		}
		close(ch)
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
			for i := 1; i < 100; i++ {
				It(strconv.Itoa(i)+" resource", func() {
					pool, err := NewWorkPoolWithMax(i)
					Expect(err).To(BeNil())
					rsp, err := pool.Run(NoOpTasks(i*2), 1, 100)
					Expect(err).To(BeNil())
					var taskResult TaskResult
					var ok bool
					for j := 0; j < i*2; j++ {
						taskResult, ok = <-rsp
						Expect(ok).To(BeTrue())
						Expect(taskResult.Result).To(Equal("NoOp"))
					}
					taskResult, ok = <-rsp
					Expect(ok).To(BeFalse())
				})
			}
		})
	})

})
