package gogen

import (
	"os/exec"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen/gomods"
)

func Get(pkg string, args ...string) *exec.Cmd {
	args = append([]string{"get"}, args...)
	args = append(args, pkg)
	cmd := exec.Command(genny.GoBin(), args...)
	return cmd
}

func Install(pkg string, args ...string) genny.RunFn {
	return func(r *genny.Runner) error {
		return gomods.Disable(func() error {
			cmd := Get(pkg, args...)
			return r.Exec(cmd)
		})
	}
}
