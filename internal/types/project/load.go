package project

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"
)

func LoadProject(path string) (*Project, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Get the real path of the file
	fullpath, err := filepath.Abs(path)
	rootDir := filepath.Dir(fullpath)
	if err != nil {
		return nil, err
	}

	project := &Project{
		Path:           rootDir,
		BuildDir:       "_build",
		DistDir:        "_dist",
		DefaultTargets: []string{"default"},
	}
	body, err := io.ReadAll(f)
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
