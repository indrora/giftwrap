/*
Copyright © 2026 Morgan Gangwere <morgan.gangwere@gmail.com>
*/
package cmd

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/indrora/giftwrap/internal/types"
	"github.com/indrora/giftwrap/internal/types/project"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
	"golang.org/x/mod/modfile"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: doInit,
}

//go:embed "default.yml"
var defaultBody []byte

func doInit(cmd *cobra.Command, args []string) {

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
	//

	pp := project.Project{
		Name: "MyProject",
		Targets: map[string]project.BuildConfig{
			"default": project.BuildConfig{
				Package: f.Module.Mod.Path,
				Targets: types.CommandList{"linux/arm64", "linux/amd64", "darwin/arm64", "darwin/amd64", "windows/arm64", "windows/amd64"},
			},
		},
	}

	// Write out the file
	o, err := os.Create(*wrapfile)
	if err != nil {
		panic(err)
	}
	defer o.Close()

	dumper, err := yaml.NewDumper(o, yaml.V4)
	if err != nil {
		panic(err)
	}
	if dumper.Dump(pp) != nil {
		fmt.Println("Couldn't write wrapfile...")
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
