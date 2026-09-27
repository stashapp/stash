package logger

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

type captureLogger struct {
	LoggerImpl
	messages []string
}

func (l *captureLogger) Debug(args ...interface{}) {
	l.messages = append(l.messages, fmt.Sprint(args[len(args)-1]))
}

func TestReadLogMessages(t *testing.T) {
	const prefix = "\x01d\x02"

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "short lines",
			input: prefix + "one\n" + prefix + "two\n",
			want:  []string{"one", "two"},
		},
		{
			name:  "no trailing newline",
			input: prefix + "one\n" + prefix + "two",
			want:  []string{"one", "two"},
		},
		{
			// a line longer than bufio.MaxScanTokenSize used to abort the read
			// entirely, discarding every subsequent line
			name:  "over-long line does not stop subsequent lines",
			input: prefix + strings.Repeat("x", maxLogLineLength*3) + "\n" + prefix + "after\n",
			want: []string{
				strings.Repeat("x", maxLogLineLength-len(prefix)) + truncatedSuffix,
				"after",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &captureLogger{}
			pl := PluginLogger{Logger: l}
			pl.ReadLogMessages(io.NopCloser(strings.NewReader(tt.input)))

			if len(l.messages) != len(tt.want) {
				t.Fatalf("got %d messages, want %d", len(l.messages), len(tt.want))
			}
			for i, want := range tt.want {
				if l.messages[i] != want {
					t.Errorf("message %d: got %.80q (len %d), want %.80q (len %d)",
						i, l.messages[i], len(l.messages[i]), want, len(want))
				}
			}
		})
	}
}
