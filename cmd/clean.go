/*
Copyright © 2026 Morgan Gangwere <morgan.gangwere@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:     "clean",
	Short:   "Clean build and release artifacts",
	Long:    `clean removes the build directory (_build) and the distribution directory (_dist).`,
	PreRunE: LoadProject,

	Run: runClean,
}

func runClean(cmd *cobra.Command, args []string) {
	rootLogger.Info("Cleaning begins")

	rootLogger.Info("Attempting to remove build directory", "dir", globProject.BuildDir)
	if err := os.RemoveAll(globProject.BuildDir); err != nil {
		rootLogger.Error("Failed to remove build directory", "dir", globProject.BuildDir, "err", err)
	} else {
		rootLogger.Info("Successfully removed build directory", "dir", globProject.BuildDir)
	}

	// And the same for the dist dir.
	rootLogger.Info("Attempting to remove dist directory", "dir", globProject.DistDir)
	if err := os.RemoveAll(globProject.DistDir); err != nil {
		rootLogger.Error("Failed to remove dist directory", "dir", globProject.DistDir, "err", err)
	} else {
		rootLogger.Info("Successfully removed dist directory", "dir", globProject.DistDir)
	}
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
