package integration

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCommands(t *testing.T) {
	tests := []struct {
		command string
		args    []string
		want    string
	}{
		{command: "01_interfaces"},
		{command: "02_methodsets"},
		{command: "03_io"},
		{command: "04_args", args: []string{"Maria", "2"}, want: "hello, Maria\nhello, Maria"},
		{command: "05_common"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.command, func(t *testing.T) {
			arguments := append([]string{"run", "../../cmd/" + tt.command}, tt.args...)
			out, err := exec.Command("go", arguments...).CombinedOutput()
			if err != nil {
				t.Fatalf("go run ./cmd/%s failed: %v\n%s", tt.command, err, out)
			}
			got := strings.TrimSpace(string(out))
			if got == "" {
				t.Fatalf("command %s returned empty output", tt.command)
			}
			if tt.want != "" && got != tt.want {
				t.Fatalf("command %s output = %q, want %q", tt.command, got, tt.want)
			}
		})
	}
}
