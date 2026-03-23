/*
Copyright © 2026 Morgan Gangwere <morgan.gangwere@gmail.com>
*/
package cmd

import (
	"github.com/indrora/giftwrap/internal/builder"
	"github.com/indrora/giftwrap/internal/packager"
	"github.com/indrora/giftwrap/internal/runner"
	"github.com/spf13/cobra"
)

var releaseCmd = &cobra.Command{
	Use:     "release",
	Short:   "Build and package project variants into release archives",
	Long:    `release builds each configured target variant and packages the output into archives in the distribution directory.`,
	Run:     doRelease,
	PreRunE: LoadProject,
}

func doRelease(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		args = globProject.DefaultTargets
	}

	run := runner.NewExecRunner(rootLogger)

	b, err := builder.NewBuilder(*globProject, run, globalRunnerOpts)
	if err != nil {
		rootLogger.Fatal("failed setting up builder", "err", err)
	}

	if releaseShell != nil && *releaseShell != "" {
		b.Shell = *releaseShell
		rootLogger.Debug("shell specified on cmdline", "shell", *releaseShell)
	}

	pkg := packager.NewPackager(*globProject)

	if err := b.Setup(); err != nil {
		rootLogger.Fatal("failed during startup", "err", err)
	}

	if err := pkg.Setup(); err != nil {
		rootLogger.Fatal("failed setting up packager", "err", err)
	}

	for _, targetName := range args {
		rootLogger.Info("releasing target", "target", targetName)

		for _, variant := range globProject.Targets[targetName].Platforms {
			rootLogger.Info("building variant", "variant", variant)

			if err := b.BuildTarget(targetName, variant); err != nil {
				rootLogger.Fatal("failed to build target", "target", targetName, "variant", variant, "err", err)
			}

			rootLogger.Info("packaging variant", "variant", variant)

			if err := pkg.PackageTarget(targetName, variant); err != nil {
				rootLogger.Fatal("failed to package target", "target", targetName, "variant", variant, "err", err)
			}
		}
	}

	if err := b.Teardown(); err != nil {
		rootLogger.Fatal("failed post-exec", "err", err)
	}

	rootLogger.Print("Finished!")
}

var releaseShell *string

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseShell = releaseCmd.Flags().String("shell", "", "Specify the shell to use for building")
}
