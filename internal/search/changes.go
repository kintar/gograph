package search

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ozgurcd/gograph/internal/graph"
)

// ChangeStatus classifies a symbol relative to the current graph.
type ChangeStatus string

const (
	// ChangeModified means the symbol's source file is newer than graph.json.
	// The symbol may or may not have changed — agents should inspect it.
	ChangeModified ChangeStatus = "modified"
	// ChangeNew means the declaration was found in a changed file but is not
	// recorded in graph.json — it was likely added since the last build.
	ChangeNew ChangeStatus = "new"
	// ChangeDeleted means a symbol from graph.json lives in a file that no
	// longer exists on disk — it was likely removed.
	ChangeDeleted ChangeStatus = "deleted"
)

// ChangedSymbol is a symbol affected by source changes since the last build.
type ChangedSymbol struct {
	// Name is the symbol name or declaration identifier.
	Name string
	// File is the source file path.
	File string
	// Line is the line number (0 for deleted symbols where the file is gone).
	Line int
	// Status classifies how this symbol was affected.
	Status ChangeStatus
}

// ChangesResult is returned by Changes.
type ChangesResult struct {
	// GraphAge is when graph.json was last generated.
	GraphAge time.Time
	// ChangedFiles lists source files newer than the graph.
	ChangedFiles []string
	// Symbols lists all symbols affected by the source changes.
	Symbols []ChangedSymbol
}

// Changes compares the current source tree against graph.json to report what
// has likely changed since the last build. It identifies:
//   - Symbols in files newer than the graph (ChangeModified)
//   - Top-level declarations in changed files not found in the graph (ChangeNew)
//   - Graph symbols whose source files no longer exist (ChangeDeleted)
//
// root is the absolute path to the repository root.
func Changes(g *graph.Graph, root string) *ChangesResult {
	graphTime := g.GeneratedAt
	result := &ChangesResult{GraphAge: graphTime}

	// Step 1: Walk source tree to find changed and deleted files.
	changedFiles := make(map[string]bool)
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			base := info.Name()
			if base == ".gograph" || base == "vendor" || base == ".git" || base == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if info.ModTime().After(graphTime) {
			if rel, err := filepath.Rel(root, path); err == nil {
				changedFiles[rel] = true
				result.ChangedFiles = append(result.ChangedFiles, rel)
			}
		}
		return nil
	})
	sortStrings(result.ChangedFiles)

	// Step 2: Build a set of all files that exist on disk (needed for deleted detection).
	existingFiles := make(map[string]bool)
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			if rel, err := filepath.Rel(root, path); err == nil {
				existingFiles[rel] = true
			}
		}
		return nil
	})

	// Step 3: Build a set of graph symbols keyed by (name, file).
	type symKey struct{ name, file string }
	graphSymbols := make(map[symKey]bool)
	for _, s := range g.Symbols {
		rel, _ := filepath.Rel(root, s.File)
		graphSymbols[symKey{s.Name, rel}] = true
		graphSymbols[symKey{s.Name, s.File}] = true
	}

	// Step 4: For each changed file, collect graph symbols (modified) and
	// parse for new declarations not in the graph.
	seenModified := make(map[symKey]bool)
	for _, s := range g.Symbols {
		rel, _ := filepath.Rel(root, s.File)
		if !changedFiles[rel] && !changedFiles[s.File] {
			continue
		}
		key := symKey{s.Name, rel}
		if seenModified[key] {
			continue
		}
		seenModified[key] = true
		result.Symbols = append(result.Symbols, ChangedSymbol{
			Name:   s.Name,
			File:   rel,
			Line:   s.Line,
			Status: ChangeModified,
		})
	}

	// Parse changed files for NEW top-level declarations.
	fset := token.NewFileSet()
	for relPath := range changedFiles {
		absPath := filepath.Join(root, relPath)
		f, err := parser.ParseFile(fset, absPath, nil, 0)
		if err != nil {
			continue
		}
		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				name := d.Name.Name
				key := symKey{name, relPath}
				if !graphSymbols[key] && !seenModified[key] {
					result.Symbols = append(result.Symbols, ChangedSymbol{
						Name:   name,
						File:   relPath,
						Line:   fset.Position(d.Pos()).Line,
						Status: ChangeNew,
					})
				}
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					switch ts := spec.(type) {
					case *ast.TypeSpec:
						name := ts.Name.Name
						key := symKey{name, relPath}
						if !graphSymbols[key] && !seenModified[key] {
							result.Symbols = append(result.Symbols, ChangedSymbol{
								Name:   name,
								File:   relPath,
								Line:   fset.Position(ts.Pos()).Line,
								Status: ChangeNew,
							})
						}
					}
				}
			}
		}
	}

	// Step 5: Detect deleted symbols — graph symbols whose files are gone.
	for _, s := range g.Symbols {
		// Check existence via os.Stat on the absolute path first.
		if _, statErr := os.Stat(s.File); statErr == nil {
			continue // file exists
		}
		// Also try relative resolution from root.
		rel, relErr := filepath.Rel(root, s.File)
		if relErr == nil && existingFiles[rel] {
			continue
		}
		result.Symbols = append(result.Symbols, ChangedSymbol{
			Name:   s.Name,
			File:   s.File,
			Line:   0,
			Status: ChangeDeleted,
		})
	}

	return result
}
