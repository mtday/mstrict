package core

import "github.com/sirupsen/logrus"

func configureLogging(debug bool) {
	formatter := &logrus.TextFormatter{
		DisableLevelTruncation:    true,
		EnvironmentOverrideColors: true,
		PadLevelText:              true,
		DisableTimestamp:          true,
		FullTimestamp:             false,
		DisableSorting:            true,
	}
	logrus.SetFormatter(formatter)

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debug logging enabled")
	}
}
