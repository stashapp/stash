package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// run_command — the assistant's escape hatch for anything the GraphQL/define_tool
// surfaces can't express: running the bundled off-queue helper scripts with custom
// flags (e.g. /usr/local/bin/identify_external.py), reading/editing config files,
// inspecting the filesystem, installing packages — i.e. changing how the app behaves,
// not just its library data.
//
// It runs `sh -c <command>` inside the Stash container and returns combined
// stdout+stderr (truncated). It is Writes:true, so under the default `ask` policy
// every command is shown to the user for approval before it runs ("with my
// permission"). Registration is gated behind `assistant_dev_loop_enabled` so the
// whole capability has an explicit on/off switch (Phase 2 gate).
//
// Boundary: it runs inside the container, so it CANNOT recompile Stash's Go binary
// (adding/altering a built-in compiled tool needs a host image rebuild). Everything
// else in the running app — config, scripts, data, files — is fair game.

const maxExecOutputBytes = 16000

// RegisterExecTools adds run_command. Caller should only invoke this when
// assistant_dev_loop_enabled is true.
func RegisterExecTools(reg *Registry) {
	reg.Register(runCommandTool())
}

func runCommandTool() *Tool {
	return &Tool{
		Name: "run_command",
		Description: "Run a shell command inside the Stash container — your escape hatch for anything the " +
			"graphql_*/define_tool surfaces can't do: run the bundled helper scripts with custom flags (e.g. " +
			"`python3 /usr/local/bin/identify_external.py --phashed-only --apply`), read or edit config files, " +
			"inspect the filesystem, install packages. Combined stdout+stderr is returned. This is a WRITE: state " +
			"exactly what the command does and why, then it runs only after the user approves. Prefer " +
			"graphql_query/graphql_mutate for library data; use this for scripts/config/filesystem. You cannot " +
			"recompile Stash itself from here.",
		Schema: json.RawMessage(`{
			"type":"object",
			"properties":{
				"command":{"type":"string","description":"Shell command, executed via 'sh -c'."},
				"timeout_seconds":{"type":"integer","description":"Kill the command after this many seconds (default 120, max 600)."}
			},
			"required":["command"]
		}`),
		Writes: true,
		Run: func(ctx context.Context, input json.RawMessage) (string, error) {
			var in struct {
				Command string `json:"command"`
				Timeout int    `json:"timeout_seconds"`
			}
			if err := json.Unmarshal(input, &in); err != nil {
				return "", fmt.Errorf("bad input: %w", err)
			}
			command := strings.TrimSpace(in.Command)
			if command == "" {
				return "", fmt.Errorf("command is required")
			}
			timeout := in.Timeout
			if timeout <= 0 {
				timeout = 120
			}
			if timeout > 600 {
				timeout = 600
			}

			cctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
			defer cancel()

			cmd := exec.CommandContext(cctx, "sh", "-c", command)
			var buf bytes.Buffer
			cmd.Stdout = &buf
			cmd.Stderr = &buf
			runErr := cmd.Run()

			out := buf.String()
			if len(out) > maxExecOutputBytes {
				out = out[:maxExecOutputBytes] + fmt.Sprintf("\n…(truncated, %d bytes total)", buf.Len())
			}
			if cctx.Err() == context.DeadlineExceeded {
				return out, fmt.Errorf("command timed out after %ds", timeout)
			}
			// Return non-zero exits as a normal result (not a tool error) so the
			// model sees the exit status + output and can react.
			if runErr != nil {
				return fmt.Sprintf("[exit: %v]\n%s", runErr, out), nil
			}
			if strings.TrimSpace(out) == "" {
				return "(exit 0, no output)", nil
			}
			return out, nil
		},
	}
}
