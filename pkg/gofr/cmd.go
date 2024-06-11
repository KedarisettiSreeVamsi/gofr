package gofr

import (
	"os"
	"regexp"
	"strings"

	cmd2 "gofr.dev/pkg/gofr/cmd"
	"gofr.dev/pkg/gofr/container"
)

type cmd struct {
	routes      []route
	description string
}

type route struct {
	pattern     string
	handler     Handler
	description string
	help        string
	fullPattern string
}

type Options func(c *route)

type ErrCommandNotFound struct{}

func (e ErrCommandNotFound) Error() string {
	return "No Command Found!" //nolint:goconst // This error is needed and repetition is in test to check for the exact string.
}

func (cmd *cmd) Validate(data []string) bool {

	for _, val := range data {
		if val != "" {
			return false
		}
	}
	return true
}

func (cmd *cmd) Run(c *container.Container) {
	args := os.Args[1:] // First one is command itself

	command := []string{}

	tempCommand := ""

	// showHelp := false

	for _, a := range args {
		if a == "" {
			continue // This takes care of cases where command has multiple spaces in between.
		}

		// if a == "-h" || a == "--help" {
		// 	showHelp = true

		// 	continue
		// }

		if a[0] != '-' {
			tempCommand = tempCommand + " " + a
		} else {
			command = append(command, tempCommand)
			tempCommand = a
		}
	}

	if tempCommand != "" {
		command = append(command, tempCommand)
	}

	// if showHelp {
	// 	cmd.printHelp()
	// }

	ctx := newContext(&cmd2.Responder{}, cmd2.NewRequest(command), c)

	for it, commandVal := range command {
		if commandVal == "" {
			continue
		}

		h := cmd.handler(commandVal)

		if h == nil {
			ctx.responder.Respond(nil, ErrCommandNotFound{})
			return
		}

		ctx.responder.Respond(h(ctx))

		if it != len(command) {
			ctx.responder.Respond("\n", nil)
		}
	}
}

func (cmd *cmd) handler(path string) Handler {
	// Trim leading dashes
	shortFlag := strings.HasPrefix(path, "-") && !strings.HasPrefix(path, "--")

	fullFlag := strings.HasPrefix(path, "--")

	if shortFlag {
		path = strings.Trim(path, "-")
	} else if fullFlag {
		path = strings.Trim(path, "--")
	}

	path = strings.Split(path, " ")[0]

	// Iterate over the routes to find a matching handler
	for _, route := range cmd.routes {

		if shortFlag {
			re := regexp.MustCompile(route.pattern)

			if cmd.Validate(re.Split(path, -1)) {
				return route.handler
			}
		}

		if fullFlag && route.fullPattern != "nil" {

			reFullPattern := regexp.MustCompile(route.fullPattern)

			if cmd.Validate(reFullPattern.Split(path, -1)) {
				return route.handler
			}
		}

	}

	// Return nil if no handler matches
	return nil
}

// AddDescription adds the description text for a specified subcommand.
func AddDescription(descString string) Options {
	return func(r *route) {
		r.description = descString
	}
}

// AddHelp adds the helper text for the given subcommand
// this is displayed when -h or --help option/flag is provided.
func AddHelp(helperString string) Options {
	return func(r *route) {
		r.help = helperString
	}
}

// AddFullPattern adds the fullPattern match for the given subcommand
// Example is --help for fullPattern and -h for pattern
func AddFullPattern(fullPattern string) Options {
	return func(r *route) {
		r.fullPattern = fullPattern
	}
}

func (cmd *cmd) addRoute(pattern string, handler Handler, options ...Options) {
	tempRoute := route{
		pattern:     pattern,
		handler:     handler,
		description: "description message not provided",
		help:        "help message not provided",
		fullPattern: "nil",
	}

	for _, opt := range options {
		opt(&tempRoute)
	}

	cmd.routes = append(cmd.routes, tempRoute)
}

// func (cmd *cmd) printHelp() {
// 	fmt.Println("Available commands:")

// 	for _, r := range cmd.routes {
// 		fmt.Printf("\n  %s\n", r.pattern)

// 		if r.description != "" {
// 			fmt.Printf("    Description: %s\n", r.description)
// 		}

// 		if r.help != "" {
// 			fmt.Printf("    Help: %s\n", r.help)
// 		}
// 	}
// }
