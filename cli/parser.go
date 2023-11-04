package cli

import "fmt"

type Parser struct {
	commands []Command
}

// Execute method request arguments of cli app without application path.
func (p *Parser) Execute(args []string) (string, error) {
	var (
		output     string
		err        error
		isExecuted = false
	)

	commandName := args[0]
	params := args[1:]

	for _, command := range p.commands {
		if commandName == command.Name {
			// Execute Command
			output, err = command.Execute(params)
			if err != nil {
				return "", err
			}

			isExecuted = true
		}
	}
	if !isExecuted {
		return "", fmt.Errorf("unknown command %s", commandName)
	}
	return output, nil
}
