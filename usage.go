package arg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// AppName is the name of the app
var AppName string

//  Version of the app
var Version string

// Usage is a short introduction of this app
var Usage string

// Author information
var Author string

// Copyright information
var Copyright string

// Fail prints usage information to stderr and exits with non-zero status
func (p *Parser) Fail(msg string) {
	p.WriteUsage(os.Stderr)
	fmt.Fprintln(os.Stderr, "error:", msg)
	os.Exit(-1)
}

// WriteUsage writes usage information to the given writer
func (p *Parser) WriteUsage(w io.Writer) {
	var positionals, options []*spec
	for _, spec := range p.spec {
		if spec.positional {
			positionals = append(positionals, spec)
		} else {
			options = append(options, spec)
		}
	}

	fmt.Fprintf(w, "usage: %s", filepath.Base(os.Args[0]))

	// write the option component of the usage message
	for _, spec := range options {
		// prefix with a space
		fmt.Fprint(w, " ")
		if !spec.required {
			fmt.Fprint(w, "[")
		}
		if spec.short != "" {
			fmt.Fprint(w, synopsis(spec, "-"+spec.short))
		} else {
			fmt.Fprint(w, synopsis(spec, "-"+spec.long))
		}
		if !spec.required {
			fmt.Fprint(w, "]")
		}
	}

	// write the positional component of the usage message
	for _, spec := range positionals {
		// prefix with a space
		fmt.Fprint(w, " ")
		up := strings.ToUpper(spec.long)
		if spec.multiple {
			fmt.Fprintf(w, "[%s [%s ...]]", up, up)
		} else {
			fmt.Fprint(w, up)
		}
	}
	fmt.Fprint(w, "\n")
}

// WriteHelp writes the usage string followed by the full help string for each option
func (p *Parser) WriteHelp(w io.Writer) {
	var positionals, options []*spec
	for _, spec := range p.spec {
		if spec.positional {
			positionals = append(positionals, spec)
		} else {
			options = append(options, spec)
		}
	}

	if AppName != "" {
		fmt.Fprintf(w, "name:\n  %s", AppName)
	}
	if Version != "" {
		fmt.Fprintf(w, " %s", Version)
	}
	if Usage != "" {
		fmt.Fprintf(w, " -- %s", Usage)
	}
	fmt.Fprint(w, "\n\n")

	if Author != "" {
		fmt.Fprintf(w, "authors:\n  %s\n\n", Author)
	}

	p.WriteUsage(w)

	// write the list of positionals
	if len(positionals) > 0 {
		fmt.Fprint(w, "\npositional arguments:\n")
		for _, spec := range positionals {
			fmt.Fprintf(w, "  %s\n", spec.long)
		}
	}

	// write the list of options
	if len(options) > 0 {
		fmt.Fprint(w, "\noptions:\n")
		const colWidth = 10
		for _, spec := range options {
			// left := "  " + synopsis(spec, "--"+spec.long)
			left := "  "
			if spec.short != "" {
				left += synopsis(spec, "-"+spec.short)
			} else {
				left += synopsis(spec, "-"+spec.long)
			}
			fmt.Fprint(w, left)
			if spec.help != "" {
				if len(left)+2 < colWidth {
					fmt.Fprint(w, strings.Repeat(" ", colWidth-len(left)))
				} else {
					fmt.Fprint(w, "\n"+strings.Repeat(" ", colWidth))
				}
				fmt.Fprint(w, spec.help)
			}
			fmt.Fprint(w, "\n")
		}
	}
	if Copyright != "" {
		fmt.Fprintf(w, "\ncopyright:\n  %s\n\n", Copyright)
	}
}

func synopsis(spec *spec, form string) string {
	if spec.dest.Kind() == reflect.Bool {
		return form
	}
	return form + " " // + strings.ToUpper(spec.long)
}
