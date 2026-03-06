package runner

import "fmt"

type ProcessFailedError struct {
	Cmd    string
	Code   int
	Reason string
}

func (e ProcessFailedError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("command %s exited with code %d: %s", e.Cmd, e.Code, e.Reason)
	}
	return fmt.Sprintf("command %s exited with code %d", e.Cmd, e.Code)
}
