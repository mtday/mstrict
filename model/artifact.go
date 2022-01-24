package model

type Artifact struct {
	Group    string `yaml:"group"    json:"group"`
	Artifact string `yaml:"artifact" json:"artifact"`
	Version  string `yaml:"version"  json:"version"`
}
