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
	return "/search <query>, /clear, /scan, /cover fetch|fetch-all, /view grid|list, /play, /open, /edit, /help, /quit"
}
