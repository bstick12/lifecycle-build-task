package task_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/bstick12/pack-lifecycle-resource/pkg/task"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"gopkg.in/src-d/go-git.v4"
)

var srcDir string

var res task.OutResponse

func populateResponse(session *gexec.Session) {
	err := json.Unmarshal(session.Out.Contents(), &res)
	Expect(err).NotTo(HaveOccurred())
}

func runCmd() (*gexec.Session, error) {

	cmd := exec.Command(bins.Task, filepath.Dir(srcDir))

	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	if err != nil {
		return &gexec.Session{}, err
	}

	return session, nil
}

var _ = Describe("Task", func() {

	var tmpDir string

	var environ []string

	BeforeEach(func() {
		var err error

		environ = os.Environ()

		tmpDir, err = ioutil.TempDir("", "out_test")
		Expect(err).ToNot(HaveOccurred())
		srcDir, err = ioutil.TempDir("", "out_src_test")
		Expect(err).ToNot(HaveOccurred())

		os.Setenv("DOCKER_CONFIG", tmpDir)

		res = task.OutResponse{}
		res.Metadata = nil
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
		Expect(os.RemoveAll(srcDir)).To(Succeed())
		task.ResetEnv(environ)
	})

	Context("invalid source", func() {

		It("fails if repository is invalid", func() {
			os.Setenv("REPOSITORY", "invalid://")
			session, err := runCmd()
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
			Expect(session.Err).Should(gbytes.Say("could not parse reference"))
		})

	})

	Context("running the lifecycle", func() {

		var session *gexec.Session

		BeforeEach(func() {
			os.Setenv("REPOSITORY", "bstick12/lifecycle-build-task-test")
			os.Setenv("TAG", "latest")
			os.Setenv("REGISTRY_USERNAME", dockerUsername)
			os.Setenv("REGISTRY_PASSWORD", dockerPassword)
			os.Setenv("CONTEXT", srcDir)
		})

		It("builds the container", func() {

			_, err := git.PlainClone(srcDir, false, &git.CloneOptions{
				URL:      "https://github.com/bstick12/goflake-server",
				Progress: os.Stdout,
			})
			Expect(err).NotTo(HaveOccurred())

			session, err = runCmd()
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 200).Should(gexec.Exit(0))

			Expect(session.Err).To(gbytes.Say("Images "))
			populateResponse(session)

			re := regexp.MustCompile(`.*Images.\((.*)\)`)
			Expect(re.Match(session.Err.Contents())).To(BeTrue())
			digest := string(re.FindSubmatch(session.Err.Contents())[1])
			Expect(digest).To(Equal(res.Version.Digest))

		})

		It("builds the container twice and uses cached layers", func() {

			_, err := git.PlainClone(srcDir, false, &git.CloneOptions{
				URL:      "https://github.com/bstick12/goflake-server",
				Progress: os.Stdout,
			})
			Expect(err).NotTo(HaveOccurred())

			cacheDir, err := ioutil.TempDir("", "cache_dir")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(cacheDir)

			os.Setenv("CACHE", cacheDir)

			session, err = runCmd()
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 200).Should(gexec.Exit(0))
			Expect(session.Err).To(gbytes.Say("Images "))
			Expect(session.Err).To(gbytes.Say("Caching layer"))

			session, err = runCmd()
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 200).Should(gexec.Exit(0))

			Expect(session.Err).To(gbytes.Say("Restoring cached layer"))
			Expect(session.Err).To(gbytes.Say("Images "))

			populateResponse(session)

			re := regexp.MustCompile(`.*Images.\((.*)\)`)
			Expect(re.Match(session.Err.Contents())).To(BeTrue())
			digest := string(re.FindSubmatch(session.Err.Contents())[1])
			Expect(digest).To(Equal(res.Version.Digest))

		})

		It("builds the container using a BP_GO_TARGETS", func() {

			os.Setenv("REPOSITORY", "bstick12/lifecycle-build-task-test-target")
			os.Setenv("BUILD_ENV", `{"BP_GO_TARGETS":"./cmd/build"}`)

			_, err := git.PlainClone(srcDir, false, &git.CloneOptions{
				URL:      "https://github.com/cloudfoundry/go-mod-cnb/",
				Progress: os.Stdout,
			})
			Expect(err).NotTo(HaveOccurred())

			session, err = runCmd()
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 200).Should(gexec.Exit(0))

			Expect(session.Err).To(gbytes.Say("Images "))

			populateResponse(session)

			re := regexp.MustCompile(`.*Images.\((.*)\)`)
			Expect(re.Match(session.Err.Contents())).To(BeTrue())
			digest := string(re.FindSubmatch(session.Err.Contents())[1])
			Expect(digest).To(Equal(res.Version.Digest))

		})

	})
})
