package sakana

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
)

// Command runs certain works, call os.Exit(4) when error occurs
type Command struct {
	*flag.FlagSet

	usageMu  sync.RWMutex
	welcome  string
	usage    example
	options  []option
	examples []example

	mu   sync.Mutex
	work func()
}

type option struct {
	names       []string
	required    bool
	description string
}

type example struct {
	usage       string
	description string
}

func NewCommand(name string, usage string, description string) (c *Command) {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	c = &Command{FlagSet: fs, usage: example{usage, description}}
	fs.Usage = c.Usage
	return
}

func (c *Command) Welcome(welcome string) {
	c.usageMu.Lock()
	defer c.usageMu.Unlock()

	c.welcome = welcome
}

func (c *Command) OptionUsage(names []string, required bool, description string) {
	c.usageMu.Lock()
	defer c.usageMu.Unlock()

	c.options = append(c.options, option{names, required, description})
}

func (c *Command) Example(usage string, description string) {
	c.usageMu.Lock()
	defer c.usageMu.Unlock()

	c.examples = append(c.examples, example{usage, description})
}

func (c *Command) Usage() {
	c.usageMu.RLock()
	defer c.usageMu.RUnlock()

	hasOption, hasExample := len(c.options) != 0, len(c.examples) != 0

	if len(c.welcome) != 0 {
		fmt.Fprintln(c.Output(), c.welcome)
	}

	fmt.Fprintf(c.Output(), "usage: %s\n", c.usage.usage)
	fmt.Fprintf(c.Output(), "    %s\n", c.usage.description)
	if hasOption || hasExample {
		fmt.Fprintln(c.Output())
	}

	if hasOption {
		fmt.Fprintln(c.Output(), "options:")
		var maxWidth int
		var existRequired bool
		for _, option := range c.options {
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
		for _, option := range c.options {
			var width int
			for i, name := range option.names {
				width += len(name)
				if i != 0 {
					fmt.Fprintf(c.Output(), ", %s", name)
				} else {
					width += 2
					fmt.Fprintf(c.Output(), "    %s", name)
				}
			}
			fmt.Fprint(c.Output(), strings.Repeat(" ", maxWidth-width+1))
			if existRequired {
				if option.required {
					fmt.Fprint(c.Output(), "(required) ")
				} else {
					fmt.Fprint(c.Output(), "           ")
				}
			}
			fmt.Fprintf(c.Output(), "- %s\n", option.description)
		}
		if hasExample {
			fmt.Fprintln(c.Output())
		}
	}

	if hasExample {
		fmt.Fprintln(c.Output(), "examples:")
		for _, example := range c.examples {
			fmt.Fprintf(c.Output(), "    %s\n", example.usage)
			fmt.Fprintf(c.Output(), "        %s\n", example.description)
		}
	}
}

// Work registers a work to invoke after parsing flags
func (c *Command) Work(work func()) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if work != nil {
		if c.work != nil {
			old := c.work
			c.work = func() { old(); work() }
		} else {
			c.work = work
		}
	}
}

func (c *Command) Serve(arguments []string) {
	c.Parse(arguments)

	if c.work != nil {
		c.work()
	}

	if c.NArg() != 0 {
		fmt.Fprintf(c.Output(), "unexpected arguments (%v)\n", c.Args())
		c.Usage()
		os.Exit(4)
	}
}
