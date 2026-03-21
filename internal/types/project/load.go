package project

import (
	"fmt"
	"io"

	"go.yaml.in/yaml/v4"
)

func LoadProject(r io.Reader, dir string) (*Project, error) {
	project := &Project{
		Path:           dir,
		BuildDir:       "_build",
		DistDir:        "_dist",
		DefaultTargets: []string{"default"},
		Exec:           &BuildCmds{},
	}
	body, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	err = yaml.Load(body, project)
	if err != nil {
		return nil, err
	}

	if len(project.Targets) < 1 {
		return nil, fmt.Errorf("No targets defined")
	}

	for _, v := range project.DefaultTargets {
		if _, exists := project.Targets[v]; !exists {
			return nil, fmt.Errorf("Default %s target does not exist", v)
		}
	}

	return project, nil
}
