package core

import (
	"flag"
)

const DefaultBuildDir = "."

type Config struct {
	Debug    bool
	BuildDir string
}

func NewConfig(parseFlags bool) *Config {
	conf := &Config{}

	// parse the command-line options
	if parseFlags {
		flag.BoolVar(&conf.Debug, "debug", false, "whether to include debug output")
		flag.StringVar(&conf.BuildDir, "build-dir", DefaultBuildDir,
			"the directory containing the source code to build, defaults to the current directory")
		flag.Parse()
	} else {
		// use suitable defaults during local testing
		conf.Debug = true
		conf.BuildDir = DefaultBuildDir
	}

	configureLogging(conf.Debug)
	return conf
}
