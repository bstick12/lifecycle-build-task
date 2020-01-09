package task

import (
	"encoding/json"
	"fmt"
	"github.com/concourse/registry-image-resource"

	"github.com/pkg/errors"
)

const DefaultTag = "latest"

type Tag string

type ConfigTask struct {
	ContextDir  string `json:"context" envconfig:"CONTEXT,optional"`
	Repository  string `json:"repository" envconfig:"REPOSITORY"`
	RawTag      Tag    `json:"tag,omitempty" envconfig:"TAG,optional"`
	Username    string `json:"username,omitempty" envconfig:"REGISTRY_USERNAME,optional"`
	Password    string `json:"password,omitempty" envconfig:"REGISTRY_PASSWORD,optional"`
	RawBuildEnv string `json:"buildenv,omitempty" envconfig:"BUILD_ENV,optional"`
	Debug       bool   `json:"debug,omitempty" envconfig:"DEBUG,optional"`
	CacheDir    string `json:"cache" envconfig:"CACHE,optional"`
}

func (config *ConfigTask) Name() string {
	return fmt.Sprintf("%s:%s", config.Repository, config.Tag())
}

func (config *ConfigTask) Tag() string {
	if config.RawTag != "" {
		return string(config.RawTag)
	}

	return DefaultTag
}

func (config *ConfigTask) BuildEnv() (map[string]string, error) {
	buildEnv := map[string]string{}
	if config.RawBuildEnv == "" {
		return buildEnv, nil
	}
	err := json.Unmarshal([]byte(config.RawBuildEnv), &buildEnv)
	if err != nil {
		return buildEnv, errors.Wrap(err, "failed to read BUILD_ENV variables")
	}
	return buildEnv, nil
}

type OutResponse struct {
	Version  resource.Version         `json:"version"`
	Metadata []resource.MetadataField `json:"metadata"`
}
