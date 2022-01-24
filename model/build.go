package model

type Build struct {
	Parent       *Build        `yaml:"-" json:"-"`
	BuildFile    string        `yaml:"-" json:"buildFile"`
	Artifact     *Artifact     `         json:"artifact"`
	Dependencies []*Dependency `         json:"dependencies"`
	Children     []*Build      `yaml:"-" json:"children"`
}
