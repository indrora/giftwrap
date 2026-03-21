/*
Copyright © 2026 Morgan Gangwere <morgan.gangwere@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

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
		b, err := getWrapfile()
		if err != nil {
			return err
		}
		*wrapfile = b

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
	"giftwrap.yml",
	".wrapfile",
	".giftwrap.yml",
	".github/giftwrap.yml",
	".github/.wrapfile",
	".github/giftwrap.yml",
}

func getWrapfile() (string, error) {
	// Check if wrapfile is empty or nil
	if wrapfile == nil || *wrapfile == "" {
		// It's empty. Look for one of the possible search strings
		for _, s := range wrapfileSearchPaths {
			_, e := os.Stat(s)
			if e == nil {
				return s, nil
			} else if !errors.Is(e, os.ErrNotExist) {
				*wrapfile = s
				return s, e
			}
		}
	} else {
		return *wrapfile, nil
	}
	return ".wrapfile", nil
}
