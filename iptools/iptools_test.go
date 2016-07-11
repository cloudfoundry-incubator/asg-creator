package iptools_test

import (
	"net"

	"github.com/cloudfoundry-incubator/asg-creator/iptools"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("iptools", func() {
	Describe("SliceNetFromRange", func() {
		var (
			ipNet  *net.IPNet
			result []iptools.IPRange
		)

		var ipRange = iptools.IPRange{
			Start: net.IP{10, 10, 0, 0},
			End:   net.IP{10, 10, 15, 255},
		}

		JustBeforeEach(func() {
			result = iptools.SliceNetFromRange(ipRange, ipNet)
		})

		Context("when the IPRange doesn't overlap the IPNet", func() {
			BeforeEach(func() {
				var err error
				_, ipNet, err = net.ParseCIDR("11.0.0.0/32")
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns a slice of one IPRange consisting of the original IPRange", func() {
				Expect(result).To(Equal([]iptools.IPRange{ipRange}))
			})
		})

		Context("when the IPRange overlaps the IPNet on the low end", func() {
			BeforeEach(func() {
				var err error
				_, ipNet, err = net.ParseCIDR("10.10.0.0/22")
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns a slice of one IPRange with the overlap removed", func() {
				Expect(result).To(Equal([]iptools.IPRange{
					{
						Start: net.IP{10, 10, 4, 0},
						End:   net.IP{10, 10, 15, 255},
					},
				}))
			})
		})

		Context("when the IPRange overlaps the IPNet on the high end", func() {
			BeforeEach(func() {
				var err error
				_, ipNet, err = net.ParseCIDR("10.10.15.0/22")
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns a slice of one IPRange with the overlap removed", func() {
				Expect(result).To(Equal([]iptools.IPRange{
					{
						Start: net.IP{10, 10, 0, 0},
						End:   net.IP{10, 10, 11, 255},
					},
				}))
			})
		})

		Context("when the IPRange overlaps the IPNet entirely", func() {
			BeforeEach(func() {
				var err error
				_, ipNet, err = net.ParseCIDR("10.10.5.0/32")
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns a slice of IPRange with the overlap removed", func() {
				Expect(result).To(Equal([]iptools.IPRange{
					{
						Start: net.IP{10, 10, 0, 0},
						End:   net.IP{10, 10, 4, 255},
					},
					{
						Start: net.IP{10, 10, 5, 1},
						End:   net.IP{10, 10, 15, 255},
					},
				}))
			})
		})

		Context("when the IPRange overlaps the IPNet exactly", func() {
			BeforeEach(func() {
				var err error
				_, ipNet, err = net.ParseCIDR("10.10.1.0/20")
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns nil", func() {
				Expect(result).To(BeNil())
			})
		})
	})
})
