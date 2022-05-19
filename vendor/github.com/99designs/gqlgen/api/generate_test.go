package api

import (
	"os"
	"path"
	"testing"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/stretchr/testify/require"
)

func cleanup(workDir string) {
	_ = os.Remove(path.Join(workDir, "server.go"))
	_ = os.RemoveAll(path.Join(workDir, "graph", "generated"))
	_ = os.Remove(path.Join(workDir, "graph", "resolver.go"))
	_ = os.Remove(path.Join(workDir, "graph", "schema.resolvers.go"))
	_ = os.Remove(path.Join(workDir, "graph", "model", "models_gen.go"))
}

func TestGenerate(t *testing.T) {
	wd, _ := os.Getwd()
	type args struct {
		workDir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				workDir: path.Join(wd, "testdata", "default"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				cleanup(tt.args.workDir)
				_ = os.Chdir(wd)
			}()
			_ = os.Chdir(tt.args.workDir)
			cfg, err := config.LoadConfigFromDefaultLocations()
			require.Nil(t, err, "failed to load config")
			if err := Generate(cfg); (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
