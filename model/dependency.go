package model

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Dependency struct {
	Group      string   `yaml:"group"      json:"group"`
	Artifact   []string `yaml:"artifact"   json:"artifact"`
	Version    string   `yaml:"version"    json:"version"`
	Type       []string `yaml:"type"       json:"type"`
	Classifier []string `yaml:"classifier" json:"classifier"`
}

func (dependency *Dependency) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.MappingNode {
		return dependency.fromMappingNode(node)
	} else if node.Kind == yaml.ScalarNode {
		return dependency.fromScalarNode(node)
	} else {
		return fmt.Errorf("unsupported dependency type on line %d", node.Line)
	}
}

func (dependency *Dependency) fromMappingNode(node *yaml.Node) error {
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]
		if key.Kind != yaml.ScalarNode {
			return fmt.Errorf("invalid dependency at line %d", key.Line)
		} else if key.Value == "group" {
			if value.Kind != yaml.ScalarNode {
				return fmt.Errorf("invalid dependency at line %d, group must be a single string value", value.Line)
			}
			dependency.Group = strings.TrimSpace(value.Value)
		} else if key.Value == "artifact" {
			if value.Kind == yaml.ScalarNode {
				dependency.Artifact = append(dependency.Artifact, strings.TrimSpace(value.Value))
			} else if value.Kind == yaml.SequenceNode {
				for _, seqValue := range value.Content {
					if seqValue.Kind != yaml.ScalarNode {
						return fmt.Errorf("invalid dependency at line %d, artifact must be a string value", value.Line)
					}
					dependency.Artifact = append(dependency.Artifact, strings.TrimSpace(seqValue.Value))
				}
			} else {
				return fmt.Errorf("invalid dependency at line %d, artifact must be one or more string values", value.Line)
			}
		} else if key.Value == "version" {
			if value.Kind != yaml.ScalarNode {
				return fmt.Errorf("invalid dependency at line %d, version must be a single string value", value.Line)
			}
			dependency.Version = strings.TrimSpace(value.Value)
		} else if key.Value == "type" {
			if value.Kind == yaml.ScalarNode {
				dependency.Type = append(dependency.Artifact, strings.TrimSpace(value.Value))
			} else if value.Kind == yaml.SequenceNode {
				for _, seqValue := range value.Content {
					if seqValue.Kind != yaml.ScalarNode {
						return fmt.Errorf("invalid dependency at line %d, type must be a string value", value.Line)
					}
					dependency.Type = append(dependency.Type, strings.TrimSpace(seqValue.Value))
				}
			} else {
				return fmt.Errorf("invalid dependency at line %d, type must be one or more string values", value.Line)
			}
		} else if key.Value == "classifier" {
			if value.Kind == yaml.ScalarNode {
				dependency.Classifier = append(dependency.Artifact, strings.TrimSpace(value.Value))
			} else if value.Kind == yaml.SequenceNode {
				for _, seqValue := range value.Content {
					if seqValue.Kind != yaml.ScalarNode {
						return fmt.Errorf("invalid dependency at line %d, classifier must be a string value", value.Line)
					}
					dependency.Classifier = append(dependency.Classifier, strings.TrimSpace(seqValue.Value))
				}
			} else {
				return fmt.Errorf("invalid dependency at line %d, classifier must be one or more string values", value.Line)
			}
		} else {
			return fmt.Errorf("unrecognized field `%s` in dependency at line %d", key.Value, key.Line)
		}
	}
	return nil
}

func (dependency *Dependency) fromScalarNode(node *yaml.Node) error {
	parts := strings.Split(node.Value, ":")
	for i, part := range parts {
		switch i {
		case 0:
			dependency.Group = strings.TrimSpace(part)
		case 1:
			for _, artifact := range strings.Split(part, ",") {
				dependency.Artifact = append(dependency.Artifact, strings.TrimSpace(artifact))
			}
		case 2:
			dependency.Version = strings.TrimSpace(part)
		case 3:
			for _, dependencyType := range strings.Split(part, ",") {
				dependency.Type = append(dependency.Type, strings.TrimSpace(dependencyType))
			}
		case 4:
			dependency.Classifier = append(dependency.Classifier, part)
			for _, classifier := range strings.Split(part, ",") {
				dependency.Classifier = append(dependency.Classifier, strings.TrimSpace(classifier))
			}
		case 5:
			return fmt.Errorf("invalid dependency, too many parts on line %d", node.Line)
		}
	}
	return nil
}
