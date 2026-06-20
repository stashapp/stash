package command

import (
	"errors"
	"strings"
)

type Command struct {
	Name string
	Args []string
	Raw  string
}

var ErrNotCommand = errors.New("input is not a slash command")

var completableCommands = []string{"search", "scan", "random", "clear", "help"}

func CompletableCommands() []string {
	return append([]string(nil), completableCommands...)
}

func Parse(input string) (Command, error) {
	input = strings.TrimSpace(input)
	if !strings.HasPrefix(input, "/") {
		return Command{}, ErrNotCommand
	}

	fields := strings.Fields(strings.TrimPrefix(input, "/"))
	if len(fields) == 0 {
		return Command{}, errors.New("slash command cannot be empty")
	}

	return Command{
		Name: strings.ToLower(fields[0]),
		Args: fields[1:],
		Raw:  input,
	}, nil
}

func Help() string {
	return ":q, /search <query>, /random <n>, /clear, /scan, /help"
}
