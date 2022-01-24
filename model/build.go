package model

type Build struct {
	Parent    *Build    `json:"-"`
	BuildFile string    `json:"buildFile"`
	Artifact  *Artifact `json:"artifact"`
	Children  []*Build  `json:"children"`
}
