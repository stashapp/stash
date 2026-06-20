package browse

import "testing"

func TestParseQuery(t *testing.T) {
	query := ParseQuery("quiet tag:demo actress:alice other")

	if query.Text != "quiet other" {
		t.Fatalf("Text = %q", query.Text)
	}
	if query.Tag != "demo" {
		t.Fatalf("Tag = %q", query.Tag)
	}
	if query.Performer != "alice" {
		t.Fatalf("Performer = %q", query.Performer)
	}
}
