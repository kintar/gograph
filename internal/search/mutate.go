package search

import (
	"strings"

	"github.com/ozgurcd/gograph/internal/graph"
)

// Mutate searches for functions that mutate the given struct field.
// The query can be "Status" or "User.Status".
func Mutate(g *graph.Graph, query string) []Result {
	parts := strings.Split(query, ".")
	field := query
	if len(parts) > 1 {
		field = parts[len(parts)-1]
	}
	field = strings.ToLower(field)

	var results []Result
	for _, m := range g.Mutations {
		if strings.ToLower(m.Field) == field {
			results = append(results, Result{
				Kind:   "mutation",
				Name:   m.Function,
				File:   m.File,
				Line:   m.Line,
				Detail: "mutates field " + m.Field,
				Score:  1,
			})
		}
	}

	sortResults(results)
	return results
}
