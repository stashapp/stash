package edit

import (
	"errors"
	"testing"
)

func TestParseArgs(t *testing.T) {
	update, err := ParseArgs([]string{"title=Demo", "rating=80", "organized=true", "watched=false"})
	if err != nil {
		t.Fatal(err)
	}

	if update.Title == nil || *update.Title != "Demo" {
		t.Fatalf("Title = %#v", update.Title)
	}
	if update.Rating == nil || *update.Rating != 80 {
		t.Fatalf("Rating = %#v", update.Rating)
	}
	if update.Organized == nil || *update.Organized != true {
		t.Fatalf("Organized = %#v", update.Organized)
	}
	if update.Watched == nil || *update.Watched != false {
		t.Fatalf("Watched = %#v", update.Watched)
	}
}

func TestBuildPartialRejectsFavorite(t *testing.T) {
	favorite := true
	_, _, err := BuildPartial(Update{Favorite: &favorite})
	if !errors.Is(err, ErrUnsupportedFavorite) {
		t.Fatalf("error = %v, want ErrUnsupportedFavorite", err)
	}
}

func TestBuildPartialValidatesRating(t *testing.T) {
	rating := 101
	_, _, err := BuildPartial(Update{Rating: &rating})
	if err == nil {
		t.Fatal("expected rating validation error")
	}
}
