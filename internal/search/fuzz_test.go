package search_test

import (
	"testing"

	"github.com/ozgurcd/gograph/internal/graph"
	"github.com/ozgurcd/gograph/internal/search"
)

// FuzzConstructors tests that the Constructors function doesn't panic
// with malformed signatures and struct names.
func FuzzConstructors(f *testing.F) {
	// Add some seed corpus
	f.Add("Worker", "func() *Worker")
	f.Add("Worker", "func(w *Worker)")
	f.Add("Server", "func() (*Server, error)")
	f.Add("", "")
	f.Add(".*", "func() .*") // regex characters

	f.Fuzz(func(t *testing.T, structName, signature string) {
		g := &graph.Graph{
			Symbols: []graph.SymbolNode{
				{
					Name:      "FuzzFunc",
					Kind:      graph.KindFunction,
					Signature: signature,
				},
			},
		}

		// Simply ensure it doesn't panic
		search.Constructors(g, structName)
	})
}

// FuzzSchema tests that Schema doesn't panic with malformed tags.
func FuzzSchema(f *testing.F) {
	// Add some seed corpus
	f.Add("users", `db:"users"`)
	f.Add("posts", `gorm:"table:posts"`)
	f.Add("", "")
	f.Add("missing", `json:"value"`)

	f.Fuzz(func(t *testing.T, tableName, tag string) {
		g := &graph.Graph{
			Symbols: []graph.SymbolNode{
				{
					Name: "FuzzStruct",
					Kind: graph.KindStruct,
					StructFields: []graph.StructField{
						{Name: "Field", Tag: tag},
					},
				},
			},
		}

		// Simply ensure it doesn't panic
		search.Schema(g, tableName)
	})
}
