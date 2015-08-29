package rancher

import (
	"strings"

	"github.com/docker/libcompose/project"
)

type SidekickInfo struct {
	primariesToSidekicks map[string][]string
	primaries            map[string]bool
	sidekickToPrimaries  map[string][]string
}

func NewSidekickInfo(project *project.Project) *SidekickInfo {
	result := &SidekickInfo{
		primariesToSidekicks: map[string][]string{},
		primaries:            map[string]bool{},
		sidekickToPrimaries:  map[string][]string{},
	}

	for name, config := range project.Configs {
		sidekicks := []string{}

		for key, value := range config.Labels.MapParts() {
			if key != "io.rancher.sidekicks" {
				continue
			}

			for _, part := range strings.Split(strings.TrimSpace(value), ",") {
				part = strings.TrimSpace(part)
				result.primaries[name] = true

				sidekicks = append(sidekicks, part)

				list, ok := result.sidekickToPrimaries[part]
				if !ok {
					list = []string{}
				}
				result.sidekickToPrimaries[part] = append(list, name)
			}
		}

		result.primariesToSidekicks[name] = sidekicks
	}

	return result
}
