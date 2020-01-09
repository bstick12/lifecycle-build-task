package task_test

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/bstick12/pack-lifecycle-resource/pkg/task"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helpers", func() {

	Context("find config directory", func() {

		It("gets config directory from DOCKER_CONFIG", func() {
			os.Setenv("DOCKER_CONFIG", "docker-config-path")
			defer os.Setenv("DOCKER_CONFIG", "")
			configDir, err := task.ConfigDir()
			Expect(err).ToNot(HaveOccurred())
			Expect(configDir).To(Equal("docker-config-path"))
		})

		It("gets default config directory", func() {
			home := os.Getenv("HOME")
			os.Setenv("HOME", "/home/user")
			defer os.Setenv("HOME", home)
			configDir, err := task.ConfigDir()
			Expect(err).ToNot(HaveOccurred())
			Expect(configDir).To(Equal("/home/user/.docker"))
		})

		It("fails if config directory can not be determined", func() {
			environ := os.Environ()
			defer task.ResetEnv(environ)
			os.Clearenv()
			_, err := task.ConfigDir()
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(task.ErrNoConfigDir))
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

			configFile, err := task.WriteConfig("registry.io", "username", "password")
			Expect(err).ToNot(HaveOccurred())
			data, err := ioutil.ReadFile(configFile)
			Expect(err).ToNot(HaveOccurred())

			var readConfig task.ConfigDocker
			err = json.Unmarshal(data, &readConfig)
			Expect(err).ToNot(HaveOccurred())
			authEntry := readConfig.Auths["registry.io"]

			decoded, err := base64.StdEncoding.DecodeString(authEntry.Auth)
			Expect(err).ToNot(HaveOccurred())

			Expect("username:password").To(Equal(string(decoded)))

		})

	})

})
