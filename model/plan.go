package model

import "github.com/mtday/mstrict/core"

type Plan struct {
	Build *Build

	ParentPlan *Plan
	SubPlans   []*Plan
}

func NewPlan(conf *core.Config, build *Build, exit bool) *Plan {
	plan := &Plan{Build: build}
	return plan
}
