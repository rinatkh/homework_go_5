package integration

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCommands(t *testing.T) {
	commands := []string{"01_interfaces", "02_methodsets", "03_io", "04_common"}
	for _, command := range commands {
		command := command
		t.Run(command, func(t *testing.T) {
			out, err := exec.Command("go", "run", "../../cmd/"+command).CombinedOutput()
			if err != nil {
				t.Fatalf("go run ./cmd/%s failed: %v\n%s", command, err, out)
			}
			if strings.TrimSpace(string(out)) == "" {
				t.Fatalf("command %s returned empty output", command)
			}
		})
	}
}
