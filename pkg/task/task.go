package task

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"

	rir "github.com/concourse/registry-image-resource"
	"github.com/fatih/color"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/otiai10/copy"
	"github.com/sirupsen/logrus"
	"github.com/vrischmann/envconfig"
)

const layersDir = "/layers"
const groupPath = layersDir + "/group.toml"
const planPath = layersDir + "/plan.toml"
const analyzedPath = layersDir + "/analyzed.toml"
const platformDir = "/platform"

type command struct {
	Cmd           string
	Flags         []string
	Writer        io.Writer
	RequiresCache bool
}

func BuildTask() {

	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	color.NoColor = false

	config := ConfigTask{}
	err := envconfig.Init(&config)
	if err != nil {
		logrus.Errorf("invalid parameters: %s", err)
		os.Exit(1)
		return
	}

	buildEnv, err := config.BuildEnv()
	if err != nil {
		logrus.Errorf("invalid parameters: %s", err)
		os.Exit(1)
		return
	}

	if config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Debugf("read config from env: %#v\n", config)

	src, err := ioutil.TempDir("", "source")
	copy.Copy(config.ContextDir, src)

	ref, err := name.ParseReference(config.Name(), name.WeakValidation)
	if err != nil {
		logrus.Errorf("could not resolve repository/tag reference: %s", err)
		os.Exit(1)
		return
	}

	cachingEnabled := false
	if config.CacheDir != "" {
		if _, err := os.Stat(config.CacheDir); err != nil {
			logrus.Errorf("cacheDir does not exist: %s", err)
		}
		cachingEnabled = true
	}

	registry := ref.Context().RegistryStr()

	configFile, err := WriteConfig(registry, config.Username, config.Password)
	if err != nil {
		logrus.Errorf("failed to write docker config.json: %s", err)
		os.Exit(1)
		return
	}
	logrus.Infof("Wrote %s for registry %s", configFile, registry)

	err = ConfigPlatformEnvVars(platformDir, buildEnv)
	if err != nil {
		logrus.Errorf("failed to write platform vars: %s", err)
		os.Exit(1)
		return
	}

	var exportBuffer bytes.Buffer
	exportWriter := io.MultiWriter(os.Stderr, &exportBuffer)

	commands := []command{
		{
			"/lifecycle/detector",
			[]string{"-app", src, "-group", groupPath, "-plan", planPath},
			os.Stderr,
			false,
		},
		{
			"/lifecycle/restorer",
			[]string{"-layers", layersDir, "-group", groupPath, "-path", config.CacheDir},
			os.Stderr,
			true,
		},
		{
			"/lifecycle/analyzer",
			[]string{"-app", src, "-layers", layersDir, "-helpers=false", "-group", groupPath, "-analyzed=" + analyzedPath, config.Name()},
			os.Stderr,
			false,
		},
		{
			"/lifecycle/builder",
			[]string{"-app", src, "-layers", layersDir, "-group", groupPath, "-plan", planPath},
			os.Stderr,
			false,
		},
		{
			"/lifecycle/exporter",
			[]string{"-app", src, "-layers", layersDir, "-helpers=false", "-group", groupPath, "-analyzed=" + analyzedPath, config.Name()},
			exportWriter,
			false,
		},
		{
			"/lifecycle/cacher",
			[]string{"-layers", layersDir, "-group", groupPath, "-path", config.CacheDir},
			os.Stderr,
			true,
		},
	}

	for _, command := range commands {
		if !command.RequiresCache || command.RequiresCache && cachingEnabled {
			cmd := exec.Command(command.Cmd, command.Flags...)
			cmd.Stdout = command.Writer
			cmd.Stderr = command.Writer
			err = cmd.Run()
			if err != nil {
				logrus.Errorf("failed to run %s: %s", command.Cmd, err)
				os.Exit(1)
				return
			}
		}
	}

	re := regexp.MustCompile(`.*Images.\((.*)\)`)
	if !re.Match(exportBuffer.Bytes()) {
		logrus.Errorf("failed to extract image digest from output")
		os.Exit(1)
		return
	}

	digest := string(re.FindSubmatch(exportBuffer.Bytes())[1])
	json.NewEncoder(os.Stdout).Encode(OutResponse{
		Version: rir.Version{
			Digest: digest,
		},
		Metadata: []rir.MetadataField{
			{
				Name:  "name",
				Value: config.Name(),
			},
		},
	})

}
