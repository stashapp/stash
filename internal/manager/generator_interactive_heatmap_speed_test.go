package manager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFunscriptDataHandlesUTF8BOM(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.funscript")
	data := append([]byte{0xEF, 0xBB, 0xBF}, []byte(`{"actions":[{"at":100,"pos":40},{"at":50,"pos":20}]}`)...)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		load func() (Script, error)
	}{
		{
			name: "load funscript data",
			load: func() (Script, error) {
				return LoadFunscriptData(path)
			},
		},
		{
			name: "heatmap generator load funscript data",
			load: func() (Script, error) {
				return NewInteractiveHeatmapSpeedGenerator(false).LoadFunscriptData(path, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			funscript, err := tt.load()
			if err != nil {
				t.Fatal(err)
			}

			if len(funscript.Actions) != 2 {
				t.Fatalf("expected 2 actions, got %d", len(funscript.Actions))
			}

			if funscript.Actions[0].At != 50 || funscript.Actions[0].Pos != 20 {
				t.Errorf("expected first action to be sorted to at=50 pos=20, got at=%v pos=%d", funscript.Actions[0].At, funscript.Actions[0].Pos)
			}
			if funscript.Actions[1].At != 100 || funscript.Actions[1].Pos != 40 {
				t.Errorf("expected second action to be sorted to at=100 pos=40, got at=%v pos=%d", funscript.Actions[1].At, funscript.Actions[1].Pos)
			}
		})
	}
}
