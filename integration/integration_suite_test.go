package integration_test

import (
	"testing"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var binPath string

var _ = SynchronizedBeforeSuite(func() []byte {
	binPath, err := gexec.Build("github.com/cloudfoundry-incubator/asg-creator")
	Expect(err).NotTo(HaveOccurred())

	return []byte(binPath)
}, func(data []byte) {
	binPath = string(data)
})

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}
