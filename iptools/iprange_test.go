package iptools_test

import (
	"net"
	"strings"

	"github.com/cloudfoundry-incubator/asg-creator/iptools"
	"github.com/cloudfoundry-incubator/candiedyaml"

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

	Describe("SliceIP", func() {
		var (
			ip     net.IP
			result []iptools.IPRange
		)

		var ipRange = iptools.IPRange{
			Start: net.IP{10, 10, 1, 0},
			End:   net.IP{10, 10, 1, 255},
		}

		JustBeforeEach(func() {
			result = ipRange.SliceIP(ip)
		})

		Context("when the IP is not in the range", func() {
			BeforeEach(func() {
				ip = net.IP{10, 10, 2, 0}
			})

			It("returns a slice of the original range", func() {
				Expect(result).To(Equal([]iptools.IPRange{ipRange}))
			})
		})

		Context("when the IP is the start of the IPRange", func() {
			BeforeEach(func() {
				ip = net.IP{10, 10, 1, 0}
			})

			It("returns a single slice", func() {
				Expect(result).To(Equal([]iptools.IPRange{
					iptools.IPRange{
						Start: net.IP{10, 10, 1, 1},
						End:   net.IP{10, 10, 1, 255},
					},
				}))
			})
		})

		Context("when the IP is the end of the IPRange", func() {
			BeforeEach(func() {
				ip = net.IP{10, 10, 1, 255}
			})

			It("returns a single slice", func() {
				Expect(result).To(Equal([]iptools.IPRange{
					iptools.IPRange{
						Start: net.IP{10, 10, 1, 0},
						End:   net.IP{10, 10, 1, 254},
					},
				}))
			})
		})

		Context("when the range only contains two IPs, and the first is being sliced out", func() {
			BeforeEach(func() {
				ipRange = iptools.IPRange{
					Start: net.IP{10, 10, 1, 0},
					End:   net.IP{10, 10, 1, 1},
				}
				ip = net.IP{10, 10, 1, 0}
			})

			XIt("returns a single slice with no End", func() {
				Expect(result).To(Equal([]iptools.IPRange{
					iptools.IPRange{
						Start: net.IP{10, 10, 1, 1},
					},
				}))
			})
		})

		Context("when the range only contains two IPs, and the last is being sliced out", func() {
			BeforeEach(func() {
				ipRange = iptools.IPRange{
					Start: net.IP{10, 10, 1, 0},
					End:   net.IP{10, 10, 1, 1},
				}
				ip = net.IP{10, 10, 1, 1}
			})

			XIt("returns a single slice with no End", func() {
				Expect(result).To(Equal([]iptools.IPRange{
					iptools.IPRange{
						Start: net.IP{10, 10, 1, 0},
					},
				}))
			})
		})

		Context("when the IP is in the IPRange", func() {
			BeforeEach(func() {
				ipRange = iptools.IPRange{
					Start: net.IP{10, 10, 1, 0},
					End:   net.IP{10, 10, 1, 255},
				}
				ip = net.IP{10, 10, 1, 5}
			})

			It("returns two slices", func() {
				Expect(result).To(Equal([]iptools.IPRange{
					iptools.IPRange{
						Start: net.IP{10, 10, 1, 0},
						End:   net.IP{10, 10, 1, 4},
					},
					iptools.IPRange{
						Start: net.IP{10, 10, 1, 6},
						End:   net.IP{10, 10, 1, 255},
					},
				}))
			})
		})
	})

	Describe("UnmarshalYAML", func() {
		var testStruct TestStruct
		var decodeErr error
		var yaml string

		JustBeforeEach(func() {
			reader := strings.NewReader(yaml)
			decoder := candiedyaml.NewDecoder(reader)
			decodeErr = decoder.Decode(&testStruct)
		})

		Context("when given valid syntax", func() {
			BeforeEach(func() {
				yaml = `
ip_range: 192.168.1.1-192.168.1.3
`
			})

			It("populates the IPRange properly", func() {
				Expect(decodeErr).NotTo(HaveOccurred())
				Expect(testStruct.IPRange.Start.Equal(net.IP{192, 168, 1, 1})).To(BeTrue())
				Expect(testStruct.IPRange.End.Equal(net.IP{192, 168, 1, 3})).To(BeTrue())
			})
		})

		Context("when given valid syntax with extra spaces", func() {
			BeforeEach(func() {
				yaml = `
ip_range: 192.168.1.1 - 192.168.1.3
`
			})

			It("populates the IPRange properly", func() {
				Expect(decodeErr).NotTo(HaveOccurred())
				Expect(testStruct.IPRange.Start.Equal(net.IP{192, 168, 1, 1})).To(BeTrue())
				Expect(testStruct.IPRange.End.Equal(net.IP{192, 168, 1, 3})).To(BeTrue())
			})
		})

		Context("when given invalid syntax", func() {
			BeforeEach(func() {
				yaml = `
ip_range: 192.168.1.1/192.168.1.3
`
			})

			It("returns an error", func() {
				Expect(decodeErr).To(HaveOccurred())
			})
		})

		Context("when given an invalid value", func() {
			BeforeEach(func() {
				yaml = `
ip_range: 192
`
			})

			It("returns an error", func() {
				Expect(decodeErr).To(HaveOccurred())
				Expect(decodeErr.Error()).To(Equal("failed-to-unmarshal-iprange-from-value: '192'"))
			})
		})
	})
})

type TestStruct struct {
	IPRange iptools.IPRange `yaml:"ip_range"`
}
