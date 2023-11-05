package cli

import "fmt"

// Parser holds a list of available commands
type Parser struct {
	commands []Command
}

// Execute finds and executes a command based on the provided arguments.
// The first argument should be the command name followed by its parameters.
// It returns the result of the command execution or an error if something goes wrong.
func (p *Parser) Execute(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("no command provided")
	}

	commandName := args[0]
	params := args[1:]

	for _, command := range p.commands {
		if commandName == command.Name {
			// Execute Command
			output, err := command.Execute(params)
			if err != nil {
				return "", err
			}
			return output, nil
		}
	}
	return "", fmt.Errorf("unknown command %s", commandName)
}
