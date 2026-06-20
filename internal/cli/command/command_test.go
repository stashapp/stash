package command

import "testing"

func TestParse(t *testing.T) {
	cmd, err := Parse("search tag:demo alice")
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

func TestParseRejectsEmptyInput(t *testing.T) {
	_, err := Parse("")
	if err == nil {
		t.Fatal("expected empty command error")
	}
}
