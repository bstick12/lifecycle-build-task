package lifecycle_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"gopkg.in/src-d/go-git.v4"

	lr "github.com/bstick12/pack-lifecycle-resource"
	resource "github.com/concourse/registry-image-resource"
)

var srcDir string

var req lr.OutRequest

var res lr.OutResponse

func populateResponse(session *gexec.Session) {
	err := json.Unmarshal(session.Out.Contents(), &res)
	Expect(err).NotTo(HaveOccurred())
}

func runCmd() (*gexec.Session, error) {

	cmd := exec.Command(bins.Out, srcDir)

	payload, err := json.Marshal(req)
	Expect(err).ToNot(HaveOccurred())

	cmd.Stdin = bytes.NewBuffer(payload)
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	if err != nil {
		return &gexec.Session{}, err
	}

	return session, nil
}

var _ = Describe("Out", func() {

	var tmpDir string

	BeforeEach(func() {
		var err error

		tmpDir, err = ioutil.TempDir("", "out_test")
		Expect(err).ToNot(HaveOccurred())
		srcDir, err = ioutil.TempDir("", "out_src_test")
		Expect(err).ToNot(HaveOccurred())

		os.Setenv("DOCKER_CONFIG", tmpDir)

		req = lr.OutRequest{}
		req.Source = resource.Source{}
		req.Params = lr.OutParams{}

		res = lr.OutResponse{}
		res.Version = resource.Version{}
		res.Metadata = nil
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
		Expect(os.RemoveAll(srcDir)).To(Succeed())
	})

	Context("invalid source", func() {

		It("fails if repository is invalid", func() {
			req.Source = resource.Source{
				Repository: "invalid://",
			}
			session, err := runCmd()
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
			Expect(session.Err).Should(gbytes.Say("could not parse reference"))
		})

	})

	Context("running the lifecycle", func() {

		var session *gexec.Session

		BeforeEach(func() {
			req.Source = resource.Source{
				Repository: "bstick12/pack-lifecycle-resource-test",
				RawTag:     "latest",
				Username:   dockerUsername,
				Password:   dockerPassword,
			}
		})

		// Docker current doesn't support delete
		// AfterEach(func() {

		// 	auth := &authn.Basic{
		// 		Username: req.Source.Username,
		// 		Password: req.Source.Password,
		// 	}

		// 	ref, err := name.ParseReference(req.Source.Name(), name.WeakValidation)
		// 	Expect(err).NotTo(HaveOccurred())

		// 	err = remote.Delete(ref, remote.WithAuth(auth))
		// 	Expect(err).NotTo(HaveOccurred())

		// })

		It("builds the container", func() {

			_, err := git.PlainClone(srcDir, false, &git.CloneOptions{
				URL:      "https://github.com/bstick12/goflake-server",
				Progress: os.Stdout,
			})
			Expect(err).NotTo(HaveOccurred())

			session, err = runCmd()
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 200).Should(gexec.Exit(0))

			Expect(session.Err).To(gbytes.Say("Resolving plan"))
			Expect(session.Err).To(gbytes.Say("Contributing app binary layer"))
			Expect(session.Err).To(gbytes.Say("Exporting layer"))
			Expect(session.Err).To(gbytes.Say("Digest:"))

			populateResponse(session)

			re := regexp.MustCompile(`.*Digest:.(.*)`)
			Expect(re.Match(session.Err.Contents())).To(BeTrue())
			digest := string(re.FindSubmatch(session.Err.Contents())[1])
			Expect(digest).To(Equal(res.Version.Digest))

		})

		XIt("builds the container twice and uses cached layers", func() {

			_, err := git.PlainClone(srcDir, false, &git.CloneOptions{
				URL:      "https://github.com/bstick12/goflake-server",
				Progress: os.Stdout,
			})
			Expect(err).NotTo(HaveOccurred())

			cacheDir, err := ioutil.TempDir("", "cache_dir")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(cacheDir)

			req.Params.CacheDir = cacheDir

			session, err = runCmd()
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 200).Should(gexec.Exit(0))
			Expect(session.Err).To(gbytes.Say("Caching layer"))
			Expect(session.Err).To(gbytes.Say("Digest:"))

			_, err = git.PlainClone(srcDir, false, &git.CloneOptions{
				URL:      "https://github.com/bstick12/goflake-server",
				Progress: os.Stdout,
			})
			Expect(err).NotTo(HaveOccurred())

			session, err = runCmd()
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 200).Should(gexec.Exit(0))

			Expect(session.Err).To(gbytes.Say("Resolving plan"))
			Expect(session.Err).To(gbytes.Say("Restoring cached layer"))
			Expect(session.Err).To(gbytes.Say("Contributing app binary layer"))
			Expect(session.Err).To(gbytes.Say("Exporting layer"))
			Expect(session.Err).To(gbytes.Say("Digest:"))

			populateResponse(session)

			re := regexp.MustCompile(`.*Digest:.(.*)`)
			Expect(re.Match(session.Err.Contents())).To(BeTrue())
			digest := string(re.FindSubmatch(session.Err.Contents())[1])
			Expect(digest).To(Equal(res.Version.Digest))

		})

	})
})
