package task_test

import (
	"encoding/base64"
	"fmt"

	"github.com/bstick12/lifecycle-build-task/pkg/task"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helpers", func() {

	Context("encodes authorization", func() {

		It("correctly for basic authentication", func() {
			encodedAuth := task.EncodeRegistryAuth("registry.io", "username", "password")
			encoded := base64.StdEncoding.EncodeToString([]byte("username:password"))
			Expect(fmt.Sprintf(`{"%s":"Basic %s"}`, "registry.io", encoded)).To(Equal(encodedAuth))
		})
	})
})
