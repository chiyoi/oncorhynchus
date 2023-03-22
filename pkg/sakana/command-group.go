package sakana

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
)

// CommandGroup routes handler by command (`FlagSet.Args()[0]`), call os.Exit(3) when error occurs
type CommandGroup struct {
	*flag.FlagSet

	usageMu  sync.RWMutex
	welcome  string
	usage    example
	options  []option
	examples []example

	mu       sync.Mutex
	commands map[string]command
	work     func()
}

type command struct {
	h     Handler
	usage string
}

func NewCommandGroup(name string, usage string, description string) *CommandGroup {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	c := &CommandGroup{FlagSet: fs, usage: example{usage, description}}
	fs.Usage = c.Usage
	return c
}

func (cg *CommandGroup) Welcome(welcome string) {
	cg.usageMu.Lock()
	defer cg.usageMu.Unlock()

	cg.welcome = welcome
}

func (cg *CommandGroup) OptionUsage(names []string, required bool, description string) {
	cg.usageMu.Lock()
	defer cg.usageMu.Unlock()

	cg.options = append(cg.options, option{names, required, description})
}

func (cg *CommandGroup) Example(usage string, description string) {
	cg.usageMu.Lock()
	defer cg.usageMu.Unlock()

	cg.examples = append(cg.examples, example{usage, description})
}

func (cg *CommandGroup) Usage() {
	cg.usageMu.RLock()
	defer cg.usageMu.RUnlock()

	hasOption, hasCommand, hasExample := len(cg.options) != 0, len(cg.commands) != 0, len(cg.examples) != 0

	if len(cg.welcome) != 0 {
		fmt.Fprintln(cg.Output(), cg.welcome)
	}

	fmt.Fprintf(cg.Output(), "usage: %s\n", cg.usage.usage)
	fmt.Fprintf(cg.Output(), "    %s\n", cg.usage.description)
	if hasOption || hasCommand || hasExample {
		fmt.Fprintln(cg.Output())
	}

	if hasOption {
		fmt.Fprintln(cg.Output(), "options:")
		var maxWidth int
		var existRequired bool
		for _, option := range cg.options {
			var width int
			for i, name := range option.names {
				width += len(name)
				if i != 1 {
					width += 2
				}
			}
			if width > maxWidth {
				maxWidth = width
			}
			existRequired = existRequired || option.required
		}
		for _, option := range cg.options {
			var width int
			for i, name := range option.names {
				width += len(name)
				if i != 0 {
					fmt.Fprintf(cg.Output(), ", %s", name)
				} else {
					width += 2
					fmt.Fprintf(cg.Output(), "    %s", name)
				}
			}
			fmt.Fprint(cg.Output(), strings.Repeat(" ", maxWidth-width+1))
			if existRequired {
				if option.required {
					fmt.Fprint(cg.Output(), "(required) ")
				} else {
					fmt.Fprint(cg.Output(), "           ")
				}
			}
			fmt.Fprintf(cg.Output(), "- %s\n", option.description)
		}
		if hasCommand || hasExample {
			fmt.Fprintln(cg.Output())
		}
	}

	if hasCommand {
		fmt.Fprintln(cg.Output(), "commands:")
		var maxWidth int
		for name := range cg.commands {
			if len(name) > maxWidth {
				maxWidth = len(name)
			}
		}
		for name, command := range cg.commands {
			fmt.Fprint(cg.Output())
			fmt.Fprintf(cg.Output(), "    %s%s - %s\n", strings.Repeat(" ", maxWidth-len(name)), name, command.usage)
		}
		if hasExample {
			fmt.Fprintln(cg.Output())
		}
	}

	if hasExample {
		fmt.Fprintln(cg.Output(), "examples:")
		for _, example := range cg.examples {
			fmt.Fprintf(cg.Output(), "    %s\n", example.usage)
			fmt.Fprintf(cg.Output(), "        %s\n", example.description)
		}
	}
}

// Command registers a command
func (cg *CommandGroup) Command(name string, h Handler, usage string) {
	cg.mu.Lock()
	defer cg.mu.Unlock()

	if cg.commands == nil {
		cg.commands = map[string]command{}
	}

	if _, ok := cg.commands[name]; ok {
		panic("duplicated command")
	}

	cg.commands[name] = command{h, usage}
}

// Work registers a work to invoke after parsing arguments
func (cg *CommandGroup) Work(work func()) {
	cg.mu.Lock()
	defer cg.mu.Unlock()

	if work != nil {
		if cg.work == nil {
			cg.work = work
		} else {
			old := cg.work
			cg.work = func() { old(); work() }
		}
	}
}

func (cg *CommandGroup) Serve(arguments []string) {
	cg.Parse(arguments)

	if cg.work != nil {
		cg.work()
	}

	if cg.NArg() <= 0 {
		fmt.Fprintln(cg.Output(), "command is required")
		cg.Usage()
		os.Exit(3)
	}

	if c, ok := cg.commands[cg.Arg(0)]; !ok {
		fmt.Fprintf(cg.Output(), "undefined command `%v`\n", cg.Arg(0))
		cg.Usage()
		os.Exit(3)
	} else {
		c.h.Serve(cg.Args()[1:])
	}
}
