package gomods

import (
	"os/exec"

	"github.com/gobuffalo/genny"
)

func Tidy(path string, verbose bool) (*genny.Generator, error) {
	g := genny.New()
	g.StepName = "go:mod:tidy:" + path
	g.RunFn(func(r *genny.Runner) error {
		if !On() {
			return nil
		}
		return r.Chdir(path, func() error {
			cmd := exec.Command(genny.GoBin(), "mod", "tidy")
			if verbose {
				cmd.Args = append(cmd.Args, "-v")
			}
			return r.Exec(cmd)
		})
	})
	return g, nil
}
