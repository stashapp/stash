package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// ErrPlugMissing ...
var ErrPlugMissing = errors.New("plugin missing")

func Decorate(c Command) *cobra.Command {
	var flags []string
	if len(c.Flags) > 0 {
		for _, f := range c.Flags {
			flags = append(flags, f)
		}
	}
	cc := &cobra.Command{
		Use:     c.Name,
		Short:   fmt.Sprintf("[PLUGIN] %s", c.Description),
		Aliases: c.Aliases,
		RunE: func(cmd *cobra.Command, args []string) error {
			plugCmd := c.Name
			if c.UseCommand != "" {
				plugCmd = c.UseCommand
			}

			ax := []string{plugCmd}
			if plugCmd == "-" {
				ax = []string{}
			}

			ax = append(ax, args...)
			ax = append(ax, flags...)

			bin, err := LookPath(c.Binary)
			if err != nil {
				return errors.WithStack(err)
			}

			ex := exec.Command(bin, ax...)
			if runtime.GOOS != "windows" {
				ex.Env = append(envy.Environ(), "BUFFALO_PLUGIN=1")
			}
			ex.Stdin = os.Stdin
			ex.Stdout = os.Stdout
			ex.Stderr = os.Stderr
			return log(strings.Join(ex.Args, " "), ex.Run)
		},
	}
	cc.DisableFlagParsing = true
	return cc
}

// LookPath ...
func LookPath(s string) (string, error) {
	if _, err := os.Stat(s); err == nil {
		return s, nil
	}

	if lp, err := exec.LookPath(s); err == nil {
		return lp, err
	}

	var bin string
	pwd, err := os.Getwd()
	if err != nil {
		return "", errors.WithStack(err)
	}

	var looks []string
	if from, err := envy.MustGet("BUFFALO_PLUGIN_PATH"); err == nil {
		looks = append(looks, from)
	} else {
		looks = []string{filepath.Join(pwd, "plugins"), filepath.Join(envy.GoPath(), "bin"), envy.Get("PATH", "")}
	}

	for _, p := range looks {
		lp := filepath.Join(p, s)
		if lp, err = filepath.EvalSymlinks(lp); err == nil {
			bin = lp
			break
		}
	}

	if len(bin) == 0 {
		return "", errors.Wrapf(ErrPlugMissing, "could not find %s in %q", s, looks)
	}
	return bin, nil
}
