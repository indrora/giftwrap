package runner

// a Runner is something that runs commands.
// How it runs commands is entirely up to the implementation
type Runner interface {
	Run(cmd string, ops Options) error
	RunArgs(cmd string, args []string, options Options) error
}
