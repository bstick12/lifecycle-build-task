package lifecycle_test

import (
	"encoding/json"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var bins struct {
	Out string `json:"out"`
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

	if _, err := os.Stat("/opt/resource/in"); err == nil {
		b.Out = "/opt/resource/out"
	} else {
		b.Out, err = gexec.Build("github.com/bstick12/pack-lifecycle-resource/cmd/out")
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
	RunSpecs(t, "PackLifecycleResource Suite")
}
