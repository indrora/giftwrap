package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/indrora/giftwrap/internal/runner"
	"github.com/indrora/giftwrap/internal/types/project"
	"github.com/spf13/cobra"
)

var globProject *project.Project

func LoadProject(cmd *cobra.Command, args []string) error {
	absPath, f, err := getWrapfile()
	if err != nil {
		return err
	}
	defer f.Close()

	dir := filepath.Dir(absPath)
	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("failed to change to project directory %s: %w", dir, err)
	}

	proj, err := project.LoadProject(f, dir)
	if err != nil {
		return err
	}

	globProject = proj
	return nil
}

var globalRunnerOpts runner.Options

func init() {

}

const (
	StdoutLevel = log.WarnLevel
	StderrLevel = log.ErrorLevel
)
