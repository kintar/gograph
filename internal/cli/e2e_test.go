package cli_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestE2E_Commands compiles a temporary binary and runs end-to-end integration tests
// to ensure the newly added commands output the expected text format.
func TestE2E_Commands(t *testing.T) {
	tmpDir := t.TempDir()

	// Write a mock package that has some constructors, globals, etc.
	mainGo := filepath.Join(tmpDir, "main.go")
	content := `package main

import "fmt"

type Worker struct {
	ID int ` + "`" + `json:"id" db:"workers"` + "`" + `
}

var GlobalCounter int

func NewWorker() *Worker {
	GlobalCounter++
	return &Worker{ID: 1}
}

type Service interface {
	DoWork()
}

func main() {
	fmt.Println("Hello", GlobalCounter)
}
`
	if err := os.WriteFile(mainGo, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write dummy source: %v", err)
	}

	// Go run cmd/gograph/main.go ... 
	// Use the compiled binary from the project root.
	repoRoot, _ := filepath.Abs("../../")
	binPath := filepath.Join(repoRoot, "bin", "gograph")
	runCmd := func(args ...string) (string, error) {
		cmd := exec.Command(binPath, args...)
		cmd.Dir = tmpDir
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		err := cmd.Run()
		return out.String(), err
	}

	// 1. Build the graph in the temporary directory
	out, err := runCmd("build", tmpDir)
	if err != nil {
		t.Fatalf("build failed: %v\nOutput:\n%s", err, out)
	}
	if !strings.Contains(out, "packages: 1") {
		t.Errorf("expected 1 package, got: %s", out)
	}

	// 2. Test constructors
	out, err = runCmd("constructors", "Worker")
	if err != nil {
		t.Fatalf("constructors failed: %v", err)
	}
	if !strings.Contains(out, "NewWorker") {
		t.Errorf("expected constructors to find NewWorker, got: %s", out)
	}

	// 3. Test schema
	out, err = runCmd("schema", "workers")
	if err != nil {
		t.Fatalf("schema failed: %v", err)
	}
	if !strings.Contains(out, "Worker") {
		t.Errorf("expected schema to find Worker, got: %s", out)
	}

	// 4. Test globals
	out, err = runCmd("globals", "main")
	if err != nil {
		t.Fatalf("globals failed: %v", err)
	}
	if !strings.Contains(out, "GlobalCounter") {
		t.Errorf("expected globals to find GlobalCounter, got: %s", out)
	}
}
