package task_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bstick12/lifecycle-build-task/pkg/task"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vrischmann/envconfig"
)

var _ = Describe("Platform Env", func() {

	var environ []string

	BeforeEach(func() {
		environ = os.Environ()
	})

	AfterEach(func() {
		task.ResetEnv(environ)
	})

	Context("Configure variables", func() {

		It("writes environment variables", func() {

			tempDir, err := ioutil.TempDir("", "platform")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tempDir)

			config := task.ConfigTask{}
			os.Setenv("REPOSITORY", "index.docker.io/username/testdata")
			os.Setenv("BUILD_ENV", `{"var1":"value1","var2":"value2"}`)

			err = envconfig.Init(&config)
			Expect(err).NotTo(HaveOccurred())
			envVars, err := config.BuildEnv()
			Expect(err).NotTo(HaveOccurred())
			err = task.ConfigPlatformEnvVars(tempDir, envVars)
			Expect(err).NotTo(HaveOccurred())

			for k, v := range envVars {
				contents, err := ioutil.ReadFile(filepath.Join(tempDir, "env", k))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(contents)).To(Equal(v))
			}

		})
	})
})
