package player

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stashapp/stash/internal/cli/browse"
)

type fakeRunner struct {
	path string
	args []string
	err  error
}

func (f *fakeRunner) Run(_ context.Context, path string, args []string) error {
	f.path = path
	f.args = append([]string(nil), args...)
	return f.err
}

type fakeRecorder struct {
	sceneID int
	called  int
}

func (f *fakeRecorder) AddView(_ context.Context, sceneID int, _ time.Time) error {
	f.sceneID = sceneID
	f.called++
	return nil
}

func TestServiceRunsFFplayAndRecordsView(t *testing.T) {
	runner := &fakeRunner{}
	recorder := &fakeRecorder{}
	service := NewWithDeps("/usr/bin/ffplay", []string{"-autoexit"}, runner, recorder)

	err := service.Play(context.Background(), browse.SceneItem{
		ID:    42,
		Title: "Scene",
		Path:  "/tmp/scene.mp4",
	})
	if err != nil {
		t.Fatal(err)
	}

	if runner.path != "/usr/bin/ffplay" {
		t.Fatalf("runner path = %q, want ffplay path", runner.path)
	}
	if want := []string{"-autoexit", "/tmp/scene.mp4"}; !reflect.DeepEqual(runner.args, want) {
		t.Fatalf("runner args = %#v, want %#v", runner.args, want)
	}
	if recorder.called != 1 || recorder.sceneID != 42 {
		t.Fatalf("recorder called=%d sceneID=%d, want one view for scene 42", recorder.called, recorder.sceneID)
	}
}

func TestServiceDoesNotRecordViewWhenFFplayFails(t *testing.T) {
	runner := &fakeRunner{err: errors.New("ffplay failed")}
	recorder := &fakeRecorder{}
	service := NewWithDeps("ffplay", nil, runner, recorder)

	err := service.Play(context.Background(), browse.SceneItem{
		ID:   42,
		Path: "/tmp/scene.mp4",
	})
	if err == nil {
		t.Fatal("expected ffplay error")
	}
	if recorder.called != 0 {
		t.Fatalf("recorder called=%d, want 0", recorder.called)
	}
}

func TestServiceRequiresVideoPath(t *testing.T) {
	service := NewWithDeps("ffplay", nil, &fakeRunner{}, &fakeRecorder{})

	err := service.Play(context.Background(), browse.SceneItem{ID: 42})
	if err == nil {
		t.Fatal("expected missing path error")
	}
}
