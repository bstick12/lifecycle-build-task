package lifecycle

import resource "github.com/concourse/registry-image-resource"

type OutParams struct {
	SourceDir string        `json:"source"`
	CacheDir  string        `json:"cache"`
	Env       []EnvVariable `json:"env"`
}

type EnvVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type OutRequest struct {
	Source resource.Source `json:"source"`
	Params OutParams       `json:"params"`
}

type OutResponse struct {
	Version  resource.Version         `json:"version"`
	Metadata []resource.MetadataField `json:"metadata"`
}
