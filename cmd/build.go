/*
Copyright © 2026 Morgan Gangwere <morgan.gangwere@gmail.com>
*/
package cmd

import (
	"github.com/indrora/giftwrap/internal/builder"
	"github.com/indrora/giftwrap/internal/runner"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:     "build",
	Short:   "Build the project",
	Long:    `build compiles each configured target variant for the specified GOOS/GOARCH pairs. If no target name is given, the default target is used.`,
	Run:     doBuild,
	PreRunE: LoadProject,
}

func doBuild(cmd *cobra.Command, args []string) {

	if len(args) == 0 {
		// use the default target

		args = globProject.DefaultTargets
	}

	run := runner.NewExecRunner(rootLogger)

	builder, err := builder.NewBuilder(*globProject, run, globalRunnerOpts)

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

		for _, a := range globProject.Targets[v].Platforms {
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

	rootLogger.Print("Finished!")

}

var shell *string

func init() {
	rootCmd.AddCommand(buildCmd)
	shell = buildCmd.Flags().String("shell", "", "Specify the shell to use for building")
}
