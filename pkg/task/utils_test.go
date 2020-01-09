package task_test

import (
	"os"

	"github.com/bstick12/lifecycle-build-task/pkg/task"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helpers", func() {

	Context("Resetting environment", func() {

		It("has the same value prior to reset", func() {
			envVars := os.Environ()
			os.Clearenv()
			Expect(os.Environ()).To(HaveLen(0))
			task.ResetEnv(envVars)
			Expect(os.Environ()).To(Equal(envVars))
		})

	})

})
