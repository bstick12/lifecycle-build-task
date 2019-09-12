package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"regexp"

	"github.com/fatih/color"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/sirupsen/logrus"

	lr "github.com/bstick12/pack-lifecycle-resource"
	resource "github.com/concourse/registry-image-resource"
)

const layersDir = "/layers"
const groupPath = layersDir + "/group.toml"
const planPath = layersDir + "/plan.toml"
const analyzedPath = layersDir + "/analyzed.toml"

type command struct {
	Cmd           string
	Flags         []string
	Writer        io.Writer
	RequiresCache bool
}

func main() {
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	color.NoColor = false

	var req lr.OutRequest
	decoder := json.NewDecoder(os.Stdin)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		logrus.Errorf("invalid payload: %s", err)
		os.Exit(1)
		return
	}

	if req.Source.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if len(os.Args) < 2 {
		logrus.Errorf("destination path not specified")
		os.Exit(1)
		return
	}

	src := os.Args[1]

	ref, err := name.ParseReference(req.Source.Name(), name.WeakValidation)
	if err != nil {
		logrus.Errorf("could not resolve repository/tag reference: %s", err)
		os.Exit(1)
		return
	}

	cachingEnabled := false
	if req.Params.CacheDir != "" {
		cachingEnabled = true
		if _, err := os.Stat(req.Params.CacheDir); err != nil {
			logrus.Errorf("cacheDir does not exist: %s", err)
		}
	}

	registry := ref.Context().RegistryStr()

	configFile, err := lr.WriteConfig(registry, req.Source.Username, req.Source.Password)
	if err != nil {
		logrus.Errorf("failed to write docker config.json: %s", err)
		os.Exit(1)
		return
	}
	logrus.Infof("Wrote %s for registry %s", configFile, registry)

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
			[]string{"-layers", layersDir, "-group", groupPath, "-path", req.Params.CacheDir},
			os.Stderr,
			true,
		},
		{
			"/lifecycle/analyzer",
			[]string{"-app", src, "-layers", layersDir, "-helpers=false", "-group", groupPath, "-analyzed=" + analyzedPath, req.Source.Name()},
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
			[]string{"-app", src, "-layers", layersDir, "-helpers=false", "-group", groupPath, "-analyzed=" + analyzedPath, req.Source.Name()},
			exportWriter,
			false,
		},
		{
			"/lifecycle/cacher",
			[]string{"-layers", layersDir, "-group", groupPath, "-path", req.Params.CacheDir},
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

	re := regexp.MustCompile(`.*Digest:.(.*)`)
	if !re.Match(exportBuffer.Bytes()) {
		logrus.Errorf("failed to extract image digest from output")
	}

	digest := string(re.FindSubmatch(exportBuffer.Bytes())[1])
	json.NewEncoder(os.Stdout).Encode(lr.OutResponse{
		Version: resource.Version{
			Digest: digest,
		},
		Metadata: []resource.MetadataField{
			{
				Name:  "name",
				Value: req.Source.Name(),
			},
		},
	})
}
