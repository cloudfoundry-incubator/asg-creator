package asg_test

import (
	"github.com/cloudfoundry-incubator/asg-creator/asg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rule", func() {
	var rule asg.Rule

	BeforeEach(func() {
		rule = asg.Rule{}
	})

	Describe("Contains", func() {
		Context("when the destination is a single IP", func() {
			BeforeEach(func() {
				rule.Destination = "127.0.0.1"
			})

			Context("when the destination matches the IP", func() {
				It("returns true", func() {
					Expect(rule.Contains("127.0.0.1")).To(BeTrue())
				})
			})

			Context("when the destination does not match the IP", func() {
				It("returns false", func() {
					Expect(rule.Contains("127.0.0.2")).To(BeFalse())
				})
			})
		})

		Context("when the destination is an IP range", func() {
			BeforeEach(func() {
				rule.Destination = "127.0.0.1-127.0.0.5"
			})

			Context("when the destination contains the IP", func() {
				It("returns true", func() {
					Expect(rule.Contains("127.0.0.1")).To(BeTrue())
				})
			})

			Context("when the destination does not contain the IP", func() {
				It("returns false", func() {
					Expect(rule.Contains("127.0.0.6")).To(BeFalse())
				})
			})
		})

		Context("when the destination is a CIDR", func() {
			BeforeEach(func() {
				rule.Destination = "127.0.0.0/24"
			})

			Context("when the destination contains the IP", func() {
				It("returns true", func() {
					Expect(rule.Contains("127.0.0.1")).To(BeTrue())
				})
			})

			Context("when the destination does not contain the IP", func() {
				It("returns false", func() {
					Expect(rule.Contains("127.0.1.0")).To(BeFalse())
				})
			})
		})
	})
})
