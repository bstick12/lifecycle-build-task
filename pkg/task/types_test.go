package task_test

import (
	"os"

	"github.com/bstick12/lifecycle-build-task/pkg/task"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vrischmann/envconfig"
)

var _ = Describe("Types", func() {

	Context("read task configuration from environment", func() {

		var environ []string

		BeforeEach(func() {
			environ = os.Environ()
			os.Setenv("REPOSITORY", "index.docker.io/username/testdata")
			os.Setenv("TAG", "v1.0")
			os.Setenv("REGISTRY_USERNAME", "username")
			os.Setenv("REGISTRY_PASSWORD", "password")
		})

		AfterEach(func() {
			task.ResetEnv(environ)
		})

		It("succesfully", func() {

			config := task.ConfigTask{}
			err := envconfig.Init(&config)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Repository).To(Equal("index.docker.io/username/testdata"))
			Expect(config.Tag()).To(Equal("v1.0"))
			Expect(config.Username).To(Equal("username"))
			Expect(config.Password).To(Equal("password"))
			Expect(config.RawBuildEnv).To(Equal(""))
			Expect(config.Debug).To(BeFalse())
		})

		It("using default tag value", func() {
			os.Setenv("TAG", "")
			config := task.ConfigTask{}
			err := envconfig.Init(&config)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Repository).To(Equal("index.docker.io/username/testdata"))
			Expect(config.Tag()).To(Equal("latest"))
			Expect(config.Username).To(Equal("username"))
			Expect(config.Password).To(Equal("password"))
			Expect(config.RawBuildEnv).To(Equal(""))
		})

		It("build environment variables", func() {
			config := task.ConfigTask{}
			os.Setenv("BUILD_ENV", `{"var1":"value1","var2":"value2"}`)
			err := envconfig.Init(&config)
			Expect(err).NotTo(HaveOccurred())
			Expect(config.RawBuildEnv).To(Equal(`{"var1":"value1","var2":"value2"}`))
			buildEnv, err := config.BuildEnv()
			Expect(buildEnv).Should(HaveKeyWithValue("var1", "value1"))
			Expect(buildEnv).Should(HaveKeyWithValue("var2", "value2"))
		})

		It("invalid build environment variables", func() {

			config := task.ConfigTask{}
			os.Setenv("BUILD_ENV", `invalid json"}`)
			err := envconfig.Init(&config)
			Expect(err).NotTo(HaveOccurred())
			_, err = config.BuildEnv()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to read BUILD_ENV"))
		})

	})

})
