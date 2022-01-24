package plan

import (
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mtday/mstrict/model"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func LoadBuilds(dir string) []*model.Build {
	directory, err := filepath.Abs(dir)
	if err != nil {
		logrus.Errorf("Failed to get absolute path for `%s`: %s", dir, err.Error())
		os.Exit(1)
	}

	dirstat, err := os.Stat(directory)
	if err != nil {
		errorMessage := err.(*os.PathError).Err.Error()
		logrus.Errorf("Failed to check build directory `%s`: %s", directory, errorMessage)
		os.Exit(1)
	}
	if !dirstat.IsDir() {
		logrus.Errorf("Not a directory: `%s`", directory)
		os.Exit(1)
	}

	var topLevelBuilds []*model.Build
	builds := map[string]*model.Build{}
	err = findBuildFiles(directory, func(path string) error {
		// find the parent build
		var parent *model.Build = nil
		var parentDir = path
		for ; len(builds) > 0 && parentDir != directory && parent == nil; {
			parentDir = filepath.Dir(parentDir)
			parent = builds[parentDir]
		}

		build := newBuild(parent, path)
		builds[filepath.Dir(path)] = build
		if parent == nil {
			topLevelBuilds = append(topLevelBuilds, build)
		} else {
			parent.Children = append(parent.Children, build)
		}
		return nil
	})
	if err != nil {
		logrus.Errorf("Failed to find build files in directory `%s`: %s", directory, err.Error())
		os.Exit(1)
	}
	return topLevelBuilds
}

func findBuildFiles(dir string, callback func(string) error) error {
	dirPtr, err := os.Open(dir)
	if err != nil {
		return err
	}
	var subdirs []string
	for {
		files, err := dirPtr.Readdir(200)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		for _, file := range files {
			if !file.IsDir() && file.Name() == model.BuildFile {
				err = callback(path.Join(dir, file.Name()))
				if err != nil {
					return err
				}
			} else if file.IsDir() && !strings.HasPrefix(file.Name(), ".") && file.Name() != model.BuildDir {
				subdirs = append(subdirs, path.Join(dir, file.Name()))
			}
		}
	}
	err = dirPtr.Close()
	if err != nil {
		return err
	}

	for _, subdir := range subdirs {
		err = findBuildFiles(subdir, callback)
		if err != nil {
			return err
		}
	}
	return nil
}

func newBuild(parent *model.Build, buildFile string) *model.Build {
	// create a build based on the provided build file
	build := &model.Build{Parent: parent, BuildFile: buildFile}

	logrus.Debugf("Reading build file: `%s`", build.BuildFile)
	file, err := os.Open(build.BuildFile)
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("Failed to close build file `%s`: %s", build.BuildFile, err.Error())
		}
	}()
	if err != nil {
		errorMessage := err.(*os.PathError).Err.Error()
		logrus.Errorf("Failed to read build file `%s`: %s", build.BuildFile, errorMessage)
		os.Exit(1)
	}

	logrus.Debugf("Parsing build file: `%s`", build.BuildFile)
	if err := yaml.NewDecoder(file).Decode(build); err != nil && err != io.EOF {
		logrus.Errorf("Failed to parse build config `%s`: %s", build.BuildFile, err.Error())
		os.Exit(1)
	}

	validateBuildArtifact(build, parent)
	return build
}

func validateBuildArtifact(build *model.Build, parent *model.Build) {
	if build.Artifact == nil {
		if parent == nil {
			logrus.Errorf("Build config `%s` is missing `artifact` field", build.BuildFile)
			os.Exit(1)
		} else {
			build.Artifact = &model.Artifact{
				Group:    parent.Artifact.Group,
				Artifact: filepath.Base(filepath.Dir(build.BuildFile)),
				Version:  parent.Artifact.Version,
			}
		}
	} else {
		if build.Artifact.Artifact == "" {
			build.Artifact.Artifact = filepath.Base(filepath.Dir(build.BuildFile))
		}
		if parent != nil {
			if build.Artifact.Group == "" {
				build.Artifact.Group = parent.Artifact.Group
			}
			if build.Artifact.Version == "" {
				build.Artifact.Version = parent.Artifact.Version
			}
		} else {
			if build.Artifact.Group == "" {
				logrus.Errorf("Build config `%s` is missing `artifact.group` field", build.BuildFile)
				os.Exit(1)
			}
			if build.Artifact.Version == "" {
				logrus.Errorf("Build config `%s` is missing `artifact.version` field", build.BuildFile)
				os.Exit(1)
			}
		}
	}

}
