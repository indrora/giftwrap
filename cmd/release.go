/*
Copyright © 2026 Morgan Gangwere <morgan.gangwere@gmail.com>
*/
package cmd

import (
	"fmt"

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

	pkg, err := packager.NewPackager(*globProject, *run)
	if err != nil {
		rootLogger.Fatal("failed setting up packager", "err", err)
	}

	if releaseShell != nil && *releaseShell != "" {
		pkg.Shell = *releaseShell
		rootLogger.Debug("shell specified on cmdline", "shell", *releaseShell)
	}

	if err := pkg.Setup(); err != nil {
		rootLogger.Fatal("failed during startup", "err", err)
	}

	for _, targetName := range args {
		rootLogger.Info("releasing target", "target", targetName)

		for _, variant := range globProject.Targets[targetName].Targets {
			rootLogger.Info("packaging variant", "variant", variant)

			if err := pkg.PackageTarget(targetName, variant); err != nil {
				rootLogger.Fatal("failed to package target", "target", targetName, "variant", variant, "err", err)
			}
		}
	}

	if err := pkg.Teardown(); err != nil {
		rootLogger.Fatal("failed post-exec", "err", err)
	}

	fmt.Println("Finished!")
}

var releaseShell *string

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseShell = releaseCmd.Flags().String("shell", "", "Specify the shell to use for building")
}
