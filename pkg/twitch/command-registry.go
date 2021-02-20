package twitch

import "fmt"

type CommandRegistry map[string]Command

type Command func() error

func (op *operator) AddCommand(name string, cmd Command) {
	op.cmdRegistry[name] = cmd
}

func (op *operator) GetCommand(name string) (Command, error) {
	cmd, ok := op.cmdRegistry[name]
	if !ok {
		return nil, fmt.Errorf("command %q is an unrecognized command", name)
	}

	return cmd, nil
}

func (op *operator) ListCommands() []string {
	cmds := []string{}

	for name := range op.cmdRegistry {
		cmds = append(cmds, name)
	}

	return cmds
}
