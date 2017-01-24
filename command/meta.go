package command

import (
	"bufio"
	"flag"
	"io"

	"github.com/hashicorp/terraform/helper/experiment"
	"github.com/hashicorp/terraform/helper/variables"

	"github.com/mitchellh/cli"
)

// Meta are the meta-options that are available on all or most commands.
type Meta struct {
	Color bool
	UI    cli.Ui

	// Variables for the context (private)
	autoKey       string
	autoVariables map[string]interface{}
	input         bool
	variables     map[string]interface{}

	// Targets for this context (private)
	targets []string

	oldUI cli.Ui
	color bool

	shadow    bool
	email     string
	password  string
	accountID int
	vOrgID    int
	create    bool
	silent    bool
	buildID   int
	name      string
}

// flags adds the meta flags to the given FlagSet.
func (m *Meta) flagSet(n string) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)
	f.BoolVar(&m.input, "input", true, "input")
	f.Var((*variables.Flag)(&m.variables), "var", "variables")
	f.Var((*variables.FlagFile)(&m.variables), "var-file", "variable file")
	f.Var((*FlagStringSlice)(&m.targets), "target", "resource to target")

	if m.autoKey != "" {
		f.Var((*variables.FlagFile)(&m.autoVariables), m.autoKey, "variable file")
	}

	// Advanced (don't need documentation, or unlikely to be set)
	f.BoolVar(&m.shadow, "shadow", true, "shadow graph")

	// Experimental features
	experiment.Flag(f)

	// Create an io.Writer that writes to our Ui properly for errors.
	// This is kind of a hack, but it does the job. Basically: create
	// a pipe, use a scanner to break it into lines, and output each line
	// to the UI. Do this forever.
	errR, errW := io.Pipe()
	errScanner := bufio.NewScanner(errR)
	go func() {
		for errScanner.Scan() {
			m.UI.Error(errScanner.Text())
		}
	}()
	f.SetOutput(errW)

	// Set the default Usage to empty
	f.Usage = func() {}

	return f
}

// process will process the meta-parameters out of the arguments. This
// will potentially modify the args in-place. It will return the resulting
// slice.
//
// vars says whether or not we support variables.
func (m *Meta) process(args []string, vars bool) []string {
	// We do this so that we retain the ability to technically call
	// process multiple times, even if we have no plans to do so
	if m.oldUI != nil {
		m.UI = m.oldUI
	}

	// Set colorization
	m.color = m.Color
	for i, v := range args {
		if v == "-no-color" {
			m.color = false
			m.Color = false
			args = append(args[:i], args[i+1:]...)
			break
		}
	}

	return args
}
