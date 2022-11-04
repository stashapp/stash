package astikit

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// SSHSession represents an SSH Session
type SSHSession interface {
	Run(string) error
	Start(string) error
	StdinPipe() (io.WriteCloser, error)
	Wait() error
}

// SSHSessionFunc represents a func that can return an SSHSession
type SSHSessionFunc func() (s SSHSession, c *Closer, err error)

// SSHCopyFileFunc is the SSH CopyFileFunc that allows doing SSH copies
func SSHCopyFileFunc(fn SSHSessionFunc) CopyFileFunc {
	return func(ctx context.Context, dst string, srcStat os.FileInfo, srcFile *os.File) (err error) {
		// Check context
		if err = ctx.Err(); err != nil {
			return
		}

		// Escape dir path
		d := strings.ReplaceAll(filepath.Dir(dst), " ", "\\ ")

		// Using local closure allows better readibility for the defer c.Close() since it
		// isolates the use of the ssh session
		if err = func() (err error) {
			// Create ssh session
			var s SSHSession
			var c *Closer
			if s, c, err = fn(); err != nil {
				err = fmt.Errorf("astikit: creating ssh session failed: %w", err)
				return
			}
			defer c.Close()

			// Create the destination folder
			if err = s.Run("mkdir -p " + d); err != nil {
				err = fmt.Errorf("astikit: creating %s failed: %w", filepath.Dir(dst), err)
				return
			}
			return
		}(); err != nil {
			return
		}

		// Using local closure allows better readibility for the defer c.Close() since it
		// isolates the use of the ssh session
		if err = func() (err error) {
			// Create ssh session
			var s SSHSession
			var c *Closer
			if s, c, err = fn(); err != nil {
				err = fmt.Errorf("astikit: creating ssh session failed: %w", err)
				return
			}
			defer c.Close()

			// Create stdin pipe
			var stdin io.WriteCloser
			if stdin, err = s.StdinPipe(); err != nil {
				err = fmt.Errorf("astikit: creating stdin pipe failed: %w", err)
				return
			}
			defer stdin.Close()

			// Use "scp" command
			if err = s.Start("scp -qt " + d); err != nil {
				err = fmt.Errorf("astikit: scp to %s failed: %w", dst, err)
				return
			}

			// Send metadata
			if _, err = fmt.Fprintln(stdin, fmt.Sprintf("C%04o", srcStat.Mode().Perm()), srcStat.Size(), filepath.Base(dst)); err != nil {
				err = fmt.Errorf("astikit: sending metadata failed: %w", err)
				return
			}

			// Copy
			if _, err = Copy(ctx, stdin, srcFile); err != nil {
				err = fmt.Errorf("astikit: copying failed: %w", err)
				return
			}

			// Send close
			if _, err = fmt.Fprint(stdin, "\x00"); err != nil {
				err = fmt.Errorf("astikit: sending close failed: %w", err)
				return
			}

			// Close stdin
			if err = stdin.Close(); err != nil {
				err = fmt.Errorf("astikit: closing failed: %w", err)
				return
			}

			// Wait
			if err = s.Wait(); err != nil {
				err = fmt.Errorf("astikit: waiting failed: %w", err)
				return
			}
			return
		}(); err != nil {
			return
		}
		return
	}
}
