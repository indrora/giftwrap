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

	run := runner.NewExecRunner(rootLogger)

	builder, err := builder.NewBuilder(*globProject, *run)

	if err != nil {
		rootLogger.Fatal("failed setting up builder", "err", err)
	}

	// Configure the shell override.

	if shell != nil && *shell != "" {
		builder.Shell = *shell
		rootLogger.Debug("Shell specified on cmdline", "shell", *shell)
	}

	rootLogger.Debug("Start build", "args", args, "dir", globProject.BuildDir)
	if err := builder.Setup(); err != nil {
		rootLogger.Fatal("failed during startup", "err", err)
	}

	for _, v := range args {
		rootLogger.Info("Building target", "target", v)

		for _, a := range globProject.Targets[v].Targets {
			rootLogger.Info("building machine", "machine", a)
			err = builder.BuildTarget(v, a)
			if err != nil {
				rootLogger.Fatal("failed to build target", "target", v, "machine", a, "err", err)
			}
		}

	}

	if err := builder.Teardown(); err != nil {
		rootLogger.Fatal("failed post-exec")
	}

	fmt.Println("Finished!")

}

var shell *string

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	shell = buildCmd.Flags().String("shell", "", "Specify the shell to use for building")
}
