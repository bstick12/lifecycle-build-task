package lifecycle

import resource "github.com/concourse/registry-image-resource"

type OutParams struct {
	CacheDir string `json:"build"`
}

type OutRequest struct {
	Source resource.Source `json:"source"`
	Params OutParams       `json:"params"`
}

type OutResponse struct {
	Version  resource.Version         `json:"version"`
	Metadata []resource.MetadataField `json:"metadata"`
}
