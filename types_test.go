package lifecycle_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	lr "github.com/bstick12/pack-lifecycle-resource"
)

var _ = Describe("Types", func() {

	Context("OutRequest", func() {
		It("parses correctly", func() {

			wd, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())
			fmt.Printf("WORKDIR = [%s]", wd)
			testFile, err := os.Open(filepath.Join("testdata", "outrequest.yml"))
			Expect(err).NotTo(HaveOccurred())
			defer testFile.Close()
			testData, err := ioutil.ReadAll(testFile)

			var req lr.OutRequest
			Expect(json.Unmarshal(testData, &req)).NotTo(HaveOccurred())

			Expect(req.Source.Repository).To(Equal("index.docker.io/username/testdata"))
			Expect(req.Source.Tag()).To(Equal("latest"))
			Expect(req.Source.Username).To(Equal("username"))
			Expect(req.Source.Password).To(Equal("password"))

			Expect(req.Params.CacheDir).To(Equal("cache_value"))
			Expect(req.Params.SourceDir).To(Equal("source_value"))
			Expect(req.Params.Env).To(HaveLen(2))

			Expect(req.Params.Env).To(ContainElement(lr.EnvVariable{Name: "ENV_1", Value: "VALUE_1"}))
			Expect(req.Params.Env).To(ContainElement(lr.EnvVariable{Name: "ENV_2", Value: "VALUE_2"}))

		})
	})
})
