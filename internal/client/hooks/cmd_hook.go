package hooks

import (
	"fmt"
	"os/exec"
)

// CommandPostHook executes given commands when a certificate has been updated.
type CommandPostHook struct {
	commands []string
}

func NewCommandPostHook(commands []string) (*CommandPostHook, error) {
	if nil == commands || len(commands) == 0 {
		return nil, fmt.Errorf("no command supplied")
	}

	return &CommandPostHook{commands: commands}, nil
}

func (hook *CommandPostHook) Invoke() error {
	cmd := exec.Command(hook.commands[0], hook.commands[1:]...)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("could not run command '%s': %v", hook.commands[0], err)
	}
	return nil
}
