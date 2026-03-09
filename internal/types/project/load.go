package project

import (
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
		Path:          rootDir,
		BuildDir:      "build",
		DistDir:       "dist",
		DefaultTarget: "default",
	}
	body, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err = yaml.Load(body, project)
	if err != nil {
		return nil, err
	}

	return project, nil
}
