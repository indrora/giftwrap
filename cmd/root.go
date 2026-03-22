/*
Copyright © 2026 Morgan Gangwere <morgan.gangwere@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/indrora/giftwrap/internal"
	"github.com/indrora/giftwrap/internal/runner"
	"github.com/spf13/cobra"
)

var rootLogger *log.Logger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "giftwrap",
	Short: "A tool to build Go applications",
	Long: `Giftwrap is a tool to build Go applications for
	multiple operating systems and architectures at a time.

	Additionally, it packages releases for you.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("No command specified. Use giftwrap init to start a project.")

		cmd.Help()

	},

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		level, err := log.ParseLevel(*loglevel)

		if err != nil {
			level = log.InfoLevel
		}

		show_caller := level == log.DebugLevel

		rootLogger = log.NewWithOptions(os.Stderr, log.Options{
			Level:           level,
			ReportTimestamp: true,
			ReportCaller:    show_caller,
		})

		globalRunnerOpts = runner.NewOptions().
			WithStdout(internal.NewLogWriter(rootLogger, StdoutLevel)).
			WithStderr(internal.NewLogWriter(rootLogger, StderrLevel))

		rootLogger.Print("G I F T W R A P !")

		return nil
	},
}

// RootCommand returns the root cobra command, for use in doc generation.
func RootCommand() *cobra.Command {
	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var wrapfile *string
var loglevel *string

func init() {
	wrapfile = rootCmd.PersistentFlags().String("wrapfile", "", "Path to the .wrapfile in use")
	loglevel = rootCmd.PersistentFlags().String("log-level", "info", "Log level to use (debug, info, warn, error)")
}

var wrapfileSearchPaths = []string{
	".wrapfile",
	"giftwrap.yml",
	".github/giftwrap.yml",
	".github/.wrapfile",
	".giftwrap.yml",
}

// getWrapfile locates the wrapfile, opens it, and returns its absolute path
// and an open file handle. The caller is responsible for closing the file.
// Returns an error if no wrapfile can be found.
func getWrapfile() (string, *os.File, error) {
	paths := wrapfileSearchPaths
	if wrapfile != nil && *wrapfile != "" {
		paths = []string{*wrapfile}
	}
	for _, p := range paths {
		f, err := os.Open(p)
		if err == nil {
			abs, err := filepath.Abs(p)
			if err != nil {
				f.Close()
				return "", nil, err
			}
			return abs, f, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return "", nil, err
		}
	}
	return "", nil, fmt.Errorf("no wrapfile found; run 'giftwrap init' to create one")
}
