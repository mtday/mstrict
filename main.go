package main

import (
	"encoding/json"

	"github.com/mtday/mstrict/core"
	"github.com/mtday/mstrict/plan"
	"github.com/sirupsen/logrus"
)

func main() {
	conf := core.NewConfig(true)

	builds := plan.LoadBuilds(conf.BuildDir)
	jsonBytes, _ := json.MarshalIndent(builds, "", "  ")
	logrus.Debugf("Builds:\n%s\n", string(jsonBytes))
}
