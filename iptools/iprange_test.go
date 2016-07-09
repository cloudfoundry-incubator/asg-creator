package iptools_test

import (
	"net"

	"github.com/cloudfoundry-incubator/asg-creator/iptools"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IPRange", func() {
	Describe("Contains", func() {
		var (
			ip     net.IP
			result bool
		)

		var ipRange = iptools.IPRange{
			Start: net.IP{10, 10, 1, 0},
			End:   net.IP{10, 10, 1, 255},
		}

		JustBeforeEach(func() {
			result = ipRange.Contains(ip)
		})

		Context("when the IP is the start of the IPRange", func() {
			BeforeEach(func() {
				ip = net.IP{10, 10, 1, 0}
			})

			It("returns true", func() {
				Expect(result).To(BeTrue())
			})
		})

		Context("when the IP is the end of the IPRange", func() {
			BeforeEach(func() {
				ip = net.IP{10, 10, 1, 255}
			})

			It("returns true", func() {
				Expect(result).To(BeTrue())
			})
		})

		Context("when the IPRange contains the IP", func() {
			BeforeEach(func() {
				ip = net.IP{10, 10, 1, 50}
			})

			It("returns true", func() {
				Expect(result).To(BeTrue())
			})
		})

		Context("when the IP is less than the IPRange", func() {
			BeforeEach(func() {
				ip = net.IP{10, 10, 0, 1}
			})

			It("returns false", func() {
				Expect(result).To(BeFalse())
			})
		})

		Context("when the IP is less than the IPRange", func() {
			BeforeEach(func() {
				ip = net.IP{10, 10, 2, 1}
			})

			It("returns false", func() {
				Expect(result).To(BeFalse())
			})
		})
	})
})
