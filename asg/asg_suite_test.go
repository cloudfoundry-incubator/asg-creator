package asg_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAsg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Asg Suite")
}
