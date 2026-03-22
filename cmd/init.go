/*
Copyright © 2026 Morgan Gangwere <morgan.gangwere@gmail.com>
*/
package cmd

import (
	_ "embed"
	"errors"
	"os"

	"github.com/charmbracelet/log"
	"github.com/indrora/giftwrap/internal/types/project"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
	"golang.org/x/mod/modfile"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project",
	Long: `Initialize a giftwrap project. This will attempt to find
a go.mod in the current file. If this does not exist, it will stop.`,
	Run: doInit,
}

func doInit(cmd *cobra.Command, args []string) {

	// Try to find the go.mod file

	_, err := os.Stat(*modpath)
	if errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Error finding go.mod: %v", err)
	}

	// if we're here, there's a go.mod in the current directory.

	data, err := os.ReadFile("go.mod")
	if err != nil {
		log.Fatalf("Error reading go.mod: %v", err)
	}

	// Parse the go.mod file
	f, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		log.Fatalf("Error parsing go.mod: %v", err)
	}

	// From this, generate a very basic configuration

	rootLogger.Debug("start building config", "modfile", f)

	pp := project.Project{
		Name: "MyProject",
		Targets: map[string]project.BuildConfig{
			"default": project.BuildConfig{
				Package: f.Module.Mod.Path,
				Platforms: []string{"linux/arm64", "linux/amd64", "darwin/arm64", "darwin/amd64", "windows/arm64", "windows/amd64"},
			},
		},
	}

	rootLogger.Debug("Generated config", "config", pp)

	// Write out the file
	outPath := *wrapfile
	if outPath == "" {
		outPath = ".wrapfile"
	}
	o, err := os.Create(outPath)
	if err != nil {
		rootLogger.Fatalf("Error creating config: %v", err)
	}
	defer o.Close()

	dumper, err := yaml.NewDumper(o, yaml.V4)
	if err != nil {
		rootLogger.Fatalf("Error creating config: %v", err)
	}
	if err = dumper.Dump(pp); err != nil {
		rootLogger.Fatalf("Error writing config: %v", err)
	}
	dumper.Close()

}

var modpath *string

func init() {
	modpath = initCmd.Flags().String("modpath", "go.mod", "Path to go.mod file")

	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
