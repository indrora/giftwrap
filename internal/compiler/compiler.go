package compiler

import (
	"encoding/json"
	"os/exec"
)

// Interface to the compiler.

// This represents a supported tool by the Go compiler
type DistTarget struct {
	GOOS         string `json:"GOOS"`
	GOARCH       string `json:"GOARCH"`
	CgoSupported bool   `json:"CgoSupported"`
	FirstClass   bool   `json:"FirstClass"`
}

// Get the list of supported targets that the Go compiler is
// willing to produce for us.
func GetDistTargets() ([]DistTarget, error) {
	targets := []DistTarget{}

	cmd := exec.Command("go", "tool", "dist", "list", "-json")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(output, &targets)

	if err != nil {
		return nil, err
	}

	return targets, nil
}
