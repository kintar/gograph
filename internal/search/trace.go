package search

import (
	"strings"

	"github.com/ozgurcd/gograph/internal/graph"
)

// TraceResult contains the error and the reverse path from entry points
type TraceResult struct {
	Error graph.ErrorEdge
	Path  []Result
}

// Trace finds the shortest path from any entry point down to the function
// that generated the given error string.
// Set includeTests to false to exclude errors and test entry points.
func Trace(g *graph.Graph, errStr string, includeTests bool) []TraceResult {
	nl := strings.ToLower(errStr)
	var matches []graph.ErrorEdge

	for _, e := range g.Errors {
		if !includeTests && isTestFile(e.File) {
			continue
		}
		if strings.Contains(strings.ToLower(e.Message), nl) {
			matches = append(matches, e)
		}
	}

	if len(matches) == 0 {
		return nil
	}

	var entryPoints []string
	for _, r := range g.Routes {
		entryPoints = append(entryPoints, r.Handler)
	}
	for _, s := range g.Symbols {
		if s.Name == "main" {
			if !includeTests && isTestFile(s.File) {
				continue
			}
			entryPoints = append(entryPoints, s.ID)
		}
	}

	var results []TraceResult
	for _, e := range matches {
		targetFunc := e.Function
		var bestPath []Result

		// Try to find the shortest path from any entry point to the error origin.
		for _, ep := range entryPoints {
			// search.Path finds the path from 'ep' to 'targetFunc'
			p := Path(g, ep, targetFunc, includeTests)
			if len(p) > 0 {
				if bestPath == nil || len(p) < len(bestPath) {
					bestPath = p
				}
			}
		}

		// If no path from an entry point, just do reverse-impact to get immediate callers
		if bestPath == nil {
			impacts := Impact(g, targetFunc, includeTests)
			if len(impacts) > 0 {
				bestPath = impacts // fallback context
			}
		}

		// (The filtering has been moved into the BFS inside Path and Impact)

		results = append(results, TraceResult{
			Error: e,
			Path:  bestPath,
		})
	}

	return results
}