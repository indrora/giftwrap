/*
Copyright © 2026 Morgan Gangwere <morgan.gangwere@gmail.com>
*/
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/indrora/giftwrap/internal"
	"github.com/indrora/giftwrap/internal/archiver"
	"github.com/indrora/giftwrap/internal/builder"
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

	singleDefault := globProject.IsSingleDefault()

	run := runner.NewExecRunner(rootLogger)

	b, err := builder.NewBuilder(*globProject, *run)
	if err != nil {
		rootLogger.Fatal("failed setting up builder", "err", err)
	}

	if releaseShell != nil && *releaseShell != "" {
		b.Shell = *releaseShell
		rootLogger.Debug("shell specified on cmdline", "shell", *releaseShell)
	}

	if err := b.Setup(); err != nil {
		rootLogger.Fatal("failed during startup", "err", err)
	}

	for _, targetName := range args {
		rootLogger.Info("releasing target", "target", targetName)

		formatStr := globProject.EffectiveArchiveFormat(targetName)
		format, err := archiver.ParseFormat(formatStr)
		if err != nil {
			rootLogger.Fatal("invalid archive format", "target", targetName, "format", formatStr, "err", err)
		}

		for _, variant := range globProject.Targets[targetName].Targets {
			rootLogger.Info("building variant", "variant", variant)

			if err := b.BuildTarget(targetName, variant); err != nil {
				rootLogger.Fatal("failed to build target", "target", targetName, "variant", variant, "err", err)
			}

			archiveName, err := archiver.ArchiveName(globProject.Name, targetName, variant, format, singleDefault)
			if err != nil {
				rootLogger.Fatal("could not construct archive name", "target", targetName, "variant", variant, "err", err)
			}

			srcDir := filepath.Join(globProject.BuildDir, internal.Slugify(targetName), internal.Slugify(variant))
			destPath := filepath.Join(globProject.DistDir, archiveName)

			rootLogger.Info("packaging", "archive", archiveName)
			if err := archiver.ArchiveDir(format, srcDir, destPath); err != nil {
				rootLogger.Fatal("failed to create archive", "archive", destPath, "err", err)
			}
		}
	}

	if err := b.Teardown(); err != nil {
		rootLogger.Fatal("failed post-exec", "err", err)
	}

	fmt.Println("Finished!")
}

var releaseShell *string

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseShell = releaseCmd.Flags().String("shell", "", "Specify the shell to use for building")
}
