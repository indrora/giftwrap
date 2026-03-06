/*
Copyright © 2026 Morgan Gangwere <morgan.gangwere@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/indrora/giftwrap/internal/builder"
	"github.com/indrora/giftwrap/internal/runner"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:     "build",
	Short:   "Build the project",
	Long:    `;)`,
	Run:     doBuild,
	PreRunE: LoadProject,
}

func doBuild(cmd *cobra.Command, args []string) {

	if len(args) == 0 {
		// use the default target

		args = []string{globProject.DefaultTarget}
	}

	run := new(runner.ExecRunner)

	builder, err := builder.NewBuilder(*globProject, *run)

	if err != nil {
		panic(err)
	}

	fmt.Println("Start build")
	if err := builder.Setup(); err != nil {
		panic(err)
	}

	for _, v := range args {
		fmt.Printf("Building target %s\n", v)
		builder.BuildTarget(v)
	}

	fmt.Println("Tearing down build")
	if err := builder.Teardown(); err != nil {
		panic(err)
	}

	fmt.Println("Finished!")

}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
