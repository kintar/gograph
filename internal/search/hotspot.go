package search

import (
	"sort"
	"strings"

	"github.com/ozgurcd/gograph/internal/graph"
)

// HotspotResult is a symbol ranked by how many call paths flow through it.
type HotspotResult struct {
	Name          string `json:"name"`
	Kind          string `json:"kind"`
	File          string `json:"file"`
	Line          int    `json:"line"`
	IncomingCalls int    `json:"incoming_calls"`
}

// Hotspot ranks all functions and methods by incoming call count (fan-in).
// The higher the count, the more central this symbol is to the codebase.
// Results are sorted descending by IncomingCalls.
// top limits the result count; pass 0 for all results.
func Hotspot(g *graph.Graph, top int) []HotspotResult {
	// Count how many times each raw callee name appears in call edges.
	incomingRaw := make(map[string]int)
	for _, c := range g.Calls {
		incomingRaw[c.CalleeRaw]++
	}

	var results []HotspotResult
	for _, s := range g.Symbols {
		if s.Kind != graph.KindFunction && s.Kind != graph.KindMethod {
			continue
		}

		displayName := s.Name
		if s.Receiver != "" {
			displayName = "(" + s.Receiver + ")." + s.Name
		}

		// Aggregate all plausible callee string forms for this symbol.
		count := 0
		count += incomingRaw[s.Name]
		count += incomingRaw[displayName]
		count += incomingRaw[s.ID]
		// Also check "pkg.FuncName" form (common for package-level calls).
		if s.PackageName != "" {
			count += incomingRaw[s.PackageName+"."+s.Name]
		}

		if count == 0 {
			continue
		}

		results = append(results, HotspotResult{
			Name:          displayName,
			Kind:          string(s.Kind),
			File:          s.File,
			Line:          s.Line,
			IncomingCalls: count,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].IncomingCalls != results[j].IncomingCalls {
			return results[i].IncomingCalls > results[j].IncomingCalls
		}
		return strings.ToLower(results[i].Name) < strings.ToLower(results[j].Name)
	})

	if top > 0 && len(results) > top {
		results = results[:top]
	}
	return results
}
