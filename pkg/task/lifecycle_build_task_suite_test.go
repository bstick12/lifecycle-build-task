package task_test

import (
	"encoding/json"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var bins struct {
	Task string `json:"task"`
}

var dockerUsername = os.Getenv("DOCKER_USERNAME")
var dockerPassword = os.Getenv("DOCKER_PASSWORD")

func checkDockerConfigured() {
	if dockerUsername == "" || dockerPassword == "" {
		Skip("must specify $DOCKER_USERNAME, and $DOCKER_PASSWORD")
	}
}

var _ = SynchronizedBeforeSuite(func() []byte {
	var err error

	checkDockerConfigured()

	b := bins

	if _, err := os.Stat("/usr/bin/lifecycle-build-task"); err == nil {
		b.Task = "/usr/bin/lifecycle-build-task"
	} else {
		b.Task, err = gexec.Build("github.com/bstick12/pack-lifecycle-resource/cmd/task")
		Expect(err).ToNot(HaveOccurred())
	}

	j, err := json.Marshal(b)
	Expect(err).ToNot(HaveOccurred())

	return j
}, func(bp []byte) {
	err := json.Unmarshal(bp, &bins)
	Expect(err).ToNot(HaveOccurred())
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	gexec.CleanupBuildArtifacts()
})

func TestPackLifecycleResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CNB Lifecycle Builder Task Suite")
}
