package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ozgurcd/gograph/internal/cli"
)

func TestBuildGraph(t *testing.T) {
	// Create a temporary directory with a dummy Go file
	tmpDir := t.TempDir()
	dummyGo := filepath.Join(tmpDir, "main.go")
	content := `package main
import "fmt"
func main() {
	fmt.Println("Hello")
}
`
	if err := os.WriteFile(dummyGo, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create dummy file: %v", err)
	}

	g, err := cli.BuildGraph(tmpDir)
	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	if g == nil {
		t.Fatal("expected non-nil graph")
	}
	if g.Root != tmpDir {
		t.Errorf("expected root %s, got %s", tmpDir, g.Root)
	}
	if len(g.Packages) == 0 {
		t.Fatal("expected at least one package")
	}
	if len(g.Files) == 0 {
		t.Fatal("expected at least one file")
	}
	if g.Files[0].PackageName != "main" {
		t.Errorf("expected package main, got %s", g.Files[0].PackageName)
	}
	
	// Check if the call was captured
	foundCall := false
	for _, call := range g.Calls {
		if call.CalleeRaw == "fmt.Println" {
			foundCall = true
		}
	}
	if !foundCall {
		t.Error("expected to find fmt.Println call in the graph")
	}
}
