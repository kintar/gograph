package search_test

import (
	"testing"

	"github.com/ozgurcd/gograph/internal/graph"
	"github.com/ozgurcd/gograph/internal/search"
)

func TestSearchFeatures(t *testing.T) {
	g := &graph.Graph{
		Symbols: []graph.SymbolNode{
			{Name: "BaseStruct", Kind: graph.KindStruct},
			{Name: "TargetStruct", Kind: graph.KindStruct, EmbeddedStructs: []string{"BaseStruct"}},
			{Name: "PublicFunc", Kind: graph.KindFunction, File: "file.go", PackageName: "test/pkg"},
			{Name: "privateFunc", Kind: graph.KindFunction, File: "file.go", PackageName: "test/pkg"},
		},
		Errors: []graph.ErrorEdge{
			{Message: "db is nil", Function: "Func1"},
			{Message: "bad request", Function: "Func2"},
		},
		SQLs: []graph.SQLEdge{
			{Query: "SELECT * FROM users", Function: "Func3"},
		},
		Packages: []graph.PackageNode{
			{
				ImportPathBestEffort: "test/pkg",
				Files:                []string{"file.go"},
			},
		},
		Files: []graph.FileNode{
			{
				Path: "file.go",
			},
		},
	}

	t.Run("Embeds", func(t *testing.T) {
		res := search.Embeds(g, "BaseStruct")
		if len(res) == 0 || res[0].Name != "TargetStruct" {
			t.Errorf("expected TargetStruct to embed BaseStruct, got %v", res)
		}
	})

	t.Run("Errors", func(t *testing.T) {
		res := search.Errors(g, "")
		if len(res) != 2 {
			t.Errorf("expected 2 error calls, got %d", len(res))
		}
	})

	t.Run("SQL", func(t *testing.T) {
		res := search.SQL(g, "SELECT")
		if len(res) != 1 || res[0].Name != "SELECT * FROM users" {
			t.Errorf("expected 1 SQL call from Func3, got %v", res)
		}
	})

	t.Run("Public", func(t *testing.T) {
		res := search.Public(g, "test/pkg")
		if len(res) != 1 || res[0].Name != "PublicFunc" {
			t.Errorf("expected 1 public symbol 'PublicFunc', got %v", res)
		}
	})
}
