package lifecycle_test

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"

	lr "github.com/bstick12/pack-lifecycle-resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helpers", func() {

	Context("find config directory", func() {

		It("gets config directory from DOCKER_CONFIG", func() {
			os.Setenv("DOCKER_CONFIG", "docker-config-path")
			defer os.Setenv("DOCKER_CONFIG", "")
			configDir, err := lr.ConfigDir()
			Expect(err).ToNot(HaveOccurred())
			Expect(configDir).To(Equal("docker-config-path"))
		})

		It("gets default config directory", func() {
			home := os.Getenv("HOME")
			os.Setenv("HOME", "/home/user")
			defer os.Setenv("HOME", home)
			configDir, err := lr.ConfigDir()
			Expect(err).ToNot(HaveOccurred())
			Expect(configDir).To(Equal("/home/user/.docker"))
		})

		It("fails if config directory can not be determined", func() {
			environ := os.Environ()
			defer lr.ResetEnv(environ)
			os.Clearenv()
			_, err := lr.ConfigDir()
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(lr.ErrNoConfigDir))
		})

	})

	Context("write config file", func() {

		It("correctly writes config file with registry and authorization", func() {
			configDir, err := ioutil.TempDir("", "docker-config")
			defer os.RemoveAll(configDir)
			Expect(err).ToNot(HaveOccurred())
			home := os.Getenv("HOME")
			os.Setenv("HOME", configDir)
			defer os.Setenv("HOME", home)

			configFile, err := lr.WriteConfig("registry.io", "username", "password")
			Expect(err).ToNot(HaveOccurred())
			data, err := ioutil.ReadFile(configFile)
			Expect(err).ToNot(HaveOccurred())

			var readConfig lr.Config
			err = json.Unmarshal(data, &readConfig)
			Expect(err).ToNot(HaveOccurred())
			authEntry := readConfig.Auths["registry.io"]

			decoded, err := base64.StdEncoding.DecodeString(authEntry.Auth)
			Expect(err).ToNot(HaveOccurred())

			Expect("username:password").To(Equal(string(decoded)))

		})

	})

	Context("Resetting environment", func() {

		It("has the same value prior to reset", func() {
			envVars := os.Environ()
			os.Clearenv()
			Expect(os.Environ()).To(HaveLen(0))
			lr.ResetEnv(envVars)
			Expect(os.Environ()).To(Equal(envVars))
		})

	})

})
