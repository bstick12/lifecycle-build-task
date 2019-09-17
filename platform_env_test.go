package lifecycle_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	lr "github.com/bstick12/pack-lifecycle-resource"
)

var _ = Describe("Platform Env", func() {

	Context("Configure variables", func() {
		It("writes environment variables", func() {

			tempDir, err := ioutil.TempDir("", "platform")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tempDir)

			envVars := []lr.EnvVariable{
				{
					Name:  "VAR1",
					Value: "VALUE1",
				},
				{
					Name:  "VAR2",
					Value: "VALUE2",
				},
			}

			for _, envVar := range envVars {
				Expect(lr.ConfigPlatformEnvVars(tempDir, envVars)).NotTo(HaveOccurred())
				contents, err := ioutil.ReadFile(filepath.Join(tempDir, "env", envVar.Name))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(contents)).To(Equal(envVar.Value))
			}

		})
	})
})
