package search

import (
	"fmt"

	"github.com/ozgurcd/gograph/internal/graph"
)

// Arity searches for functions and methods that have an argument count (arity)
// greater than or equal to the minArgs threshold. It helps find "Long Parameter List"
// code smells which are prime candidates for refactoring into Options structs.
func Arity(g *graph.Graph, minArgs int) []Result {
	var results []Result

	for _, s := range g.Symbols {
		if s.Kind == graph.KindFunction || s.Kind == graph.KindMethod {
			if s.Arity >= minArgs {
				results = append(results, Result{
					Kind:   "arity",
					Name:   s.Name,
					File:   s.File,
					Line:   s.Line,
					Detail: fmt.Sprintf("%d arguments", s.Arity),
					Score:  s.Arity, // Higher arity scores higher so it sorts to the top
				})
			}
		}
	}

	sortResults(results)
	return results
}
