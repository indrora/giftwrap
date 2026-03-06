package cmd

import (
	"github.com/indrora/giftwrap/internal/types/project"
	"github.com/spf13/cobra"
)

var globProject *project.Project

func LoadProject(cmd *cobra.Command, args []string) error {
	// Load the project from disk

	projectPath, err := getWrapfile()
	if err != nil {
		return err
	}

	// load it

	proj, err := project.LoadProject(projectPath)
	if err != nil {
		return err
	}

	globProject = proj

	return nil
}
