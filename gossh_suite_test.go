package gossh_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGossh(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gossh Suite")
}
