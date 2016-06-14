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

		Context("when the config contains public_networks: true", func() {
			BeforeEach(func() {
				config = `
public_networks: true
`
			})

			AfterEach(func() {
				os.RemoveAll("public-networks.json")
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
		})

		Context("when the config contains private_networks: true", func() {
			BeforeEach(func() {
				config = `
private_networks: true
`
			})

			AfterEach(func() {
				os.RemoveAll("private-networks.json")
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

			Context("when the config contains private networks to blacklist", func() {
				BeforeEach(func() {
					config = `
private_networks: true
excluded_networks:
- 192.168.1.0/24
`
				})

				It("should not include an ASG for that network", func() {
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

			Context("when the config contains IPs to blacklist", func() {
				BeforeEach(func() {
					config = `
private_networks: true
excluded_ips:
- 192.168.100.4
`
				})

				It("should create an ASG that skips that IP", func() {
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
    "destination": "192.168.100.5-192.168.255.255"
  }
]`)))
				})
			})

			Context("when the config contains both networks and IPs to blacklist", func() {
				BeforeEach(func() {
					config = `
private_networks: true
excluded_ips:
- 192.168.100.4
excluded_networks:
- 192.168.1.0/24
`
				})

				It("should create an ASG that skips both", func() {
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
		})
	})

	Context("when not given a configFile", func() {
		BeforeEach(func() {
			cmd = exec.Command(binPath, "create")
		})

		It("exits with error", func() {
			sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(sess).Should(gexec.Exit(1))
		})
	})
})
