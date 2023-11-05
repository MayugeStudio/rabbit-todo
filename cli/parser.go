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

// AddCommand adds a new command to the parser.
// It checks for duplicate command names to avoid conflicts.
// If a command with the same name already exists, it returns an error.
// Otherwise, it appends the new command to the parser's list of commands.
func (p *Parser) AddCommand(command Command) error {
	for _, c := range p.commands {
		if c.Name == command.Name {
			return fmt.Errorf("duplicate command name %s", command.Name)
		}
	}
	p.commands = append(p.commands, command)
	return nil
}
