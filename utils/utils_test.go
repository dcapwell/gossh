package utils_test

import (
	. "github.com/dcapwell/gossh/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
  "strconv"
  "log"
)

var _ = Describe("Utils", func() {
  Context("expand hostnames", func() {
    Context("with no expand params", func() {
      It("should return input", func() {
        input := "helloworld"
        output, err := Expand(input)
        Expect(err).To(BeNil())
        Expect(output).To(HaveLen(1))
        Expect(output[0]).To(Equal(input))
      })
    })
    Context("with single range expand param", func() {
      It("should return multiple outputs", func() {
        input := "example-[1..3]"
        output, err := Expand(input)
        Expect(err).To(BeNil())
        Expect(output).To(HaveLen(3))
        for i := 0; i < 3; i++ {
          Expect(output[i]).To(Equal("example-" + strconv.Itoa(i + 1)))
        }
      })
    })
    Context("with multiple range expand params", func() {
      It("should return multiple outputs", func() {
        input := "example-[1..3]-[1..3]"
        output, err := Expand(input)
        log.Printf("Multiple Range Expand Result: %v\n", output)
        Expect(err).To(BeNil())
        Expect(output).To(HaveLen(9))
        count := 0
        for i := 0; i < 3; i++ {
          base := "example-" + strconv.Itoa(i + 1) + "-"
          for j := 0; j < 3; j++ {
            Expect(output[count]).To(Equal(base + strconv.Itoa(j + 1)))
            count++
          }
        }
      })
    })
    Context("with single list param", func() {
      It("should return mulitple outputs", func() {
        input := "hello,world"
        output, err := Expand(input)
        Expect(err).To(BeNil())
        Expect(output).To(HaveLen(2))
        Expect(output[0]).To(Equal("hello"))
        Expect(output[1]).To(Equal("world"))
      })
    })
    Context("with multiple list params", func() {
      It("should return mulitple outputs", func() {
        input := "hello,world,how,are,you"
        output, err := Expand(input)
        Expect(err).To(BeNil())
        Expect(output).To(HaveLen(5))
        Expect(output[0]).To(Equal("hello"))
        Expect(output[1]).To(Equal("world"))
        Expect(output[2]).To(Equal("how"))
        Expect(output[3]).To(Equal("are"))
        Expect(output[4]).To(Equal("you"))
      })
    })
    Context("with single list and single range expand  params", func() {
      It("should return multiple outputs", func() {
        input := "example-[1..3],helloworld"
        output, err := Expand(input)
        Expect(err).To(BeNil())
        Expect(output).To(HaveLen(4))
        for i := 0; i < 3; i++ {
          Expect(output[i]).To(Equal("example-" + strconv.Itoa(i + 1)))
        }
        Expect(output[3]).To(Equal("helloworld"))
      })
    })
    Context("with single list and multiple range expand  params", func() {
      It("should return multiple outputs", func() {
        input := "example-[1..3]-[1..3],helloworld"
        output, err := Expand(input)
        Expect(err).To(BeNil())
        Expect(output).To(HaveLen(10))
        count := 0
        for i := 0; i < 3; i++ {
          base := "example-" + strconv.Itoa(i + 1) + "-"
          for j := 0; j < 3; j++ {
            Expect(output[count]).To(Equal(base + strconv.Itoa(j + 1)))
            count++
          }
        }
        Expect(output[9]).To(Equal("helloworld"))
      })
    })
    Context("with multiple list and single range expand  params", func() {
      It("should return multiple outputs", func() {
        input := "FIRST,example-[1..3],helloworld"
        output, err := Expand(input)
        Expect(err).To(BeNil())
        Expect(output).To(HaveLen(5))
        Expect(output[0]).To(Equal("FIRST"))
        for i := 0; i < 3; i++ {
          Expect(output[i + 1]).To(Equal("example-" + strconv.Itoa(i + 1)))
        }
        Expect(output[4]).To(Equal("helloworld"))
      })
    })
    Context("with multiple list and multiple range expand  params", func() {
      It("should return multiple outputs", func() {
        input := "FIRST,example-[1..3]-[1..3],helloworld"
        output, err := Expand(input)
        Expect(err).To(BeNil())
        Expect(output).To(HaveLen(11))
        Expect(output[0]).To(Equal("FIRST"))
        count := 1
        for i := 0; i < 3; i++ {
          base := "example-" + strconv.Itoa(i + 1) + "-"
          for j := 0; j < 3; j++ {
            Expect(output[count]).To(Equal(base + strconv.Itoa(j + 1)))
            count++
          }
        }
        Expect(output[10]).To(Equal("helloworld"))
      })
    })
  })
  Context("with invalid range", func() {
    /*
    It("should return error", func() {
      //TODO this should be valid.  Its logical to want your range to go backwards.
      _ := "invalid-[10..1]"
    })
    */
    It("non int be elipse should fail", func() {
    })
    It("multiple elipse should fail", func() {
    })
  })
})
