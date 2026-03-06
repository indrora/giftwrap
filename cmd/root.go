/*
Copyright © 2026 Morgan Gangwere <morgan.gangwere@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		b, err := getWrapfile()
		if err != nil {
			panic(err)
		}

		*wrapfile = b

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

func init() {
	wrapfile = rootCmd.PersistentFlags().String("wrapfile", "", "Path to the .wrapfile in use")
}

var wrapfileSearchPaths = []string{
	".wrapfile",
	"giftwrap.yml",
	".github/giftwrap.yml",
	".github/.wrapfile",
	".giftwrap.yml",
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
