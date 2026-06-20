package command

import (
	"errors"
	"testing"
)

func TestParse(t *testing.T) {
	cmd, err := Parse("/search tag:demo alice")
	if err != nil {
		t.Fatal(err)
	}
	if cmd.Name != "search" {
		t.Fatalf("Name = %q", cmd.Name)
	}
	if len(cmd.Args) != 2 || cmd.Args[0] != "tag:demo" || cmd.Args[1] != "alice" {
		t.Fatalf("Args = %#v", cmd.Args)
	}
}

func TestParseRejectsPlainInput(t *testing.T) {
	_, err := Parse("search")
	if !errors.Is(err, ErrNotCommand) {
		t.Fatalf("error = %v, want ErrNotCommand", err)
	}
}
