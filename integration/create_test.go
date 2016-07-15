package integration_test

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", func() {
	var cmd *exec.Cmd

	AfterEach(func() {
		os.RemoveAll("public-networks.json")
		os.RemoveAll("private-networks.json")
	})

	Context("when not given a config", func() {
		BeforeEach(func() {
			cmd = exec.Command(binPath, "create")
		})

		It("writes public-networks.json", func() {
			sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(sess).Should(gexec.Exit(0))

			_, err = os.Lstat("public-networks.json")
			Expect(err).NotTo(HaveOccurred())

			bs, err := ioutil.ReadFile("public-networks.json")
			Expect(err).NotTo(HaveOccurred())

			Expect(bs).To(MatchJSON([]byte(`
				[
					{
						"destination": "0.0.0.0-9.255.255.255",
						"protocol": "all"
					},
					{
						"destination": "11.0.0.0-169.254.169.253",
						"protocol": "all"
					},
					{
						"destination": "169.254.169.255-172.15.255.255",
						"protocol": "all"
					},
					{
						"destination": "172.32.0.0-192.167.255.255",
						"protocol": "all"
					},
					{
						"destination": "192.169.0.0-255.255.255.255",
						"protocol": "all"
					}
				]`)))
		})

		It("writes private-networks.json", func() {
			sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(sess).Should(gexec.Exit(0))

			_, err = os.Lstat("private-networks.json")
			Expect(err).NotTo(HaveOccurred())

			bs, err := ioutil.ReadFile("private-networks.json")
			Expect(err).NotTo(HaveOccurred())

			Expect(bs).To(MatchJSON([]byte(`
				[
					{
						"protocol": "all",
						"destination": "10.0.0.0-10.255.255.255"
					},
					{
						"protocol": "all",
						"destination": "172.16.0.0-172.31.255.255"
					},
					{
						"protocol": "all",
						"destination": "192.168.0.0-192.168.255.255"
					}
				]`)))
		})

	})

	Context("when given a config", func() {
		var configFile *os.File
		var config string

		JustBeforeEach(func() {
			var err error
			configFile, err = ioutil.TempFile("", "")
			Expect(err).NotTo(HaveOccurred())

			err = ioutil.WriteFile(configFile.Name(), []byte(config), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			cmd = exec.Command(binPath, "create", "--config", configFile.Name())
		})

		AfterEach(func() {
			os.RemoveAll(configFile.Name())
		})

		Context("when the config contains private networks to exclude", func() {
			BeforeEach(func() {
				config = `
excluded_networks:
- 192.168.1.0/24
`
			})

			It("should omit that network in the private-networks ASG", func() {
				sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess).Should(gexec.Exit(0))

				_, err = os.Lstat("private-networks.json")
				Expect(err).NotTo(HaveOccurred())

				bs, err := ioutil.ReadFile("private-networks.json")
				Expect(err).NotTo(HaveOccurred())

				Expect(bs).To(MatchJSON([]byte(`
					[
						{
							"protocol": "all",
							"destination": "10.0.0.0-10.255.255.255"
						},
						{
							"protocol": "all",
							"destination": "172.16.0.0-172.31.255.255"
						},
						{
							"protocol": "all",
							"destination": "192.168.0.0-192.168.0.255"
						},
						{
							"protocol": "all",
							"destination": "192.168.2.0-192.168.255.255"
						}
					]`)))
			})
		})

		Context("when the config contains public networks to exclude", func() {
			BeforeEach(func() {
				config = `
excluded_networks:
- 11.0.1.0/24
`
			})

			It("should omit that network in the public-networks ASG", func() {
				sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess).Should(gexec.Exit(0))

				_, err = os.Lstat("public-networks.json")
				Expect(err).NotTo(HaveOccurred())

				bs, err := ioutil.ReadFile("public-networks.json")
				Expect(err).NotTo(HaveOccurred())

				Expect(bs).To(MatchJSON([]byte(`
					[
						{
							"protocol": "all",
							"destination": "0.0.0.0-9.255.255.255"
						},
						{
							"protocol": "all",
							"destination": "11.0.0.0-11.0.0.255"
						},
						{
							"protocol": "all",
							"destination": "11.0.2.0-169.254.169.253"
						},
						{
							"protocol": "all",
							"destination": "169.254.169.255-172.15.255.255"
						},
						{
							"protocol": "all",
							"destination": "172.32.0.0-192.167.255.255"
						},
						{
							"protocol": "all",
							"destination": "192.169.0.0-255.255.255.255"
						}
					]`)))
			})
		})

		Context("when the config contains private IPs to exclude", func() {
			BeforeEach(func() {
				config = `
excluded_ips:
- 192.168.100.4
- 192.168.100.8
`
			})

			It("should create an ASG that skips those IPs", func() {
				sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess).Should(gexec.Exit(0))

				_, err = os.Lstat("private-networks.json")
				Expect(err).NotTo(HaveOccurred())

				bs, err := ioutil.ReadFile("private-networks.json")
				Expect(err).NotTo(HaveOccurred())

				Expect(bs).To(MatchJSON([]byte(`
					[
						{
							"protocol": "all",
							"destination": "10.0.0.0-10.255.255.255"
						},
						{
							"protocol": "all",
							"destination": "172.16.0.0-172.31.255.255"
						},
						{
							"protocol": "all",
							"destination": "192.168.0.0-192.168.100.3"
						},
						{
							"protocol": "all",
							"destination": "192.168.100.5-192.168.100.7"
						},
						{
							"protocol": "all",
							"destination": "192.168.100.9-192.168.255.255"
						}
					]`)))
			})
		})

		Context("when the config contains public IPs to exclude", func() {
			BeforeEach(func() {
				config = `
excluded_ips:
- 11.0.0.5
- 11.0.0.8
`
			})

			It("should create an ASG that skips those IPs", func() {
				sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess).Should(gexec.Exit(0))

				_, err = os.Lstat("public-networks.json")
				Expect(err).NotTo(HaveOccurred())

				bs, err := ioutil.ReadFile("public-networks.json")
				Expect(err).NotTo(HaveOccurred())

				Expect(bs).To(MatchJSON([]byte(`
				  [
				  	{
				  		"destination": "0.0.0.0-9.255.255.255",
				  		"protocol": "all"
				  	},
				  	{
				  		"destination": "11.0.0.0-11.0.0.4",
				  		"protocol": "all"
				  	},
				  	{
				  		"destination": "11.0.0.6-11.0.0.7",
				  		"protocol": "all"
				  	},
				  	{
				  		"destination": "11.0.0.9-169.254.169.253",
				  		"protocol": "all"
				  	},
				  	{
				  		"destination": "169.254.169.255-172.15.255.255",
				  		"protocol": "all"
				  	},
				  	{
				  		"destination": "172.32.0.0-192.167.255.255",
				  		"protocol": "all"
				  	},
				  	{
				  		"destination": "192.169.0.0-255.255.255.255",
				  		"protocol": "all"
				  	}
				  ]`)))
			})
		})

		Context("when the config contains both networks and IPs to exclude", func() {
			BeforeEach(func() {
				config = `
excluded_ips:
- 192.168.100.4
excluded_networks:
- 192.168.1.0/24
`
			})

			It("should create an ASG that omits both", func() {
				sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess).Should(gexec.Exit(0))

				_, err = os.Lstat("private-networks.json")
				Expect(err).NotTo(HaveOccurred())

				bs, err := ioutil.ReadFile("private-networks.json")
				Expect(err).NotTo(HaveOccurred())

				Expect(bs).To(MatchJSON([]byte(`
					[
						{
							"protocol": "all",
							"destination": "10.0.0.0-10.255.255.255"
						},
						{
							"protocol": "all",
							"destination": "172.16.0.0-172.31.255.255"
						},
						{
							"protocol": "all",
							"destination": "192.168.0.0-192.168.0.255"
						},
						{
							"protocol": "all",
							"destination": "192.168.2.0-192.168.100.3"
						},
						{
							"protocol": "all",
							"destination": "192.168.100.5-192.168.255.255"
						}
					]`)))
			})
		})

		Context("when the config is such that it expects a rule with a single IP", func() {
			BeforeEach(func() {
				config = `
excluded_ips:
- 192.168.100.4
- 192.168.100.6
`
			})

			It("should create a rule with a single IP", func() {
				sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess).Should(gexec.Exit(0))

				_, err = os.Lstat("private-networks.json")
				Expect(err).NotTo(HaveOccurred())

				bs, err := ioutil.ReadFile("private-networks.json")
				Expect(err).NotTo(HaveOccurred())

				Expect(bs).To(MatchJSON([]byte(`
					[
						{
							"protocol": "all",
							"destination": "10.0.0.0-10.255.255.255"
						},
						{
							"protocol": "all",
							"destination": "172.16.0.0-172.31.255.255"
						},
						{
							"protocol": "all",
							"destination": "192.168.0.0-192.168.100.3"
						},
						{
							"protocol": "all",
							"destination": "192.168.100.5"
						},
						{
							"protocol": "all",
							"destination": "192.168.100.7-192.168.255.255"
						}
					]`)))
			})
		})
	})

	Context("when given a config and an output file", func() {
		var configFile *os.File
		var outputFile *os.File
		var config string

		JustBeforeEach(func() {
			var err error
			configFile, err = ioutil.TempFile("", "")
			Expect(err).NotTo(HaveOccurred())

			err = ioutil.WriteFile(configFile.Name(), []byte(config), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			outputFile, err = ioutil.TempFile("", "")
			Expect(err).NotTo(HaveOccurred())

			cmd = exec.Command(binPath, "create", "--config", configFile.Name(), "--output", outputFile.Name())
		})

		AfterEach(func() {
			os.RemoveAll(configFile.Name())
			os.RemoveAll(outputFile.Name())
		})

		Context("when the config contains networks to include", func() {
			BeforeEach(func() {
				config = `
included_networks:
- 10.68.192.0/24

excluded_ips:
- 10.68.192.0
- 10.68.192.127
- 10.68.192.128
- 10.68.192.255
`
			})

			It("should create rules that only include that network", func() {
				sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess).Should(gexec.Exit(0))

				_, err = os.Lstat(outputFile.Name())
				Expect(err).NotTo(HaveOccurred())

				bs, err := ioutil.ReadFile(outputFile.Name())
				Expect(err).NotTo(HaveOccurred())

				Expect(bs).To(MatchJSON([]byte(`[
					{
							"protocol": "all",
							"destination": "10.68.192.1-10.68.192.126"
					},
					{
							"protocol": "all",
							"destination": "10.68.192.129-10.68.192.254"
					}
				]`)))
			})
		})

		Context("when the config contains IP ranges to exclude", func() {
			BeforeEach(func() {
				config = `
included_networks:
- 10.68.192.0/24

excluded_ranges:
- 10.68.192.0-10.68.192.5
- 10.68.192.127-10.68.192.128
`
			})

			It("should create rules that exclude those networks", func() {
				sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess).Should(gexec.Exit(0))

				_, err = os.Lstat(outputFile.Name())
				Expect(err).NotTo(HaveOccurred())

				bs, err := ioutil.ReadFile(outputFile.Name())
				Expect(err).NotTo(HaveOccurred())

				Expect(bs).To(MatchJSON([]byte(`[
					{
							"protocol": "all",
							"destination": "10.68.192.6-10.68.192.126"
					},
					{
							"protocol": "all",
							"destination": "10.68.192.129-10.68.192.255"
					}
				]`)))
			})
		})
	})
})
