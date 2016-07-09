package iptools_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIptools(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Iptools Suite")
}
