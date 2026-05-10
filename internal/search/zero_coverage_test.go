package search_test

import (
	"strings"
	"testing"

	"github.com/ozgurcd/gograph/internal/graph"
	"github.com/ozgurcd/gograph/internal/search"
)

func buildCoverageGraph() *graph.Graph {
	return &graph.Graph{
		Packages: []graph.PackageNode{
			{Name: "api", Dir: "pkg/api", Files: []string{"pkg/api/server.go"}},
		},
		Symbols: []graph.SymbolNode{
			{Name: "Server", Kind: graph.KindStruct, PackageName: "api", File: "pkg/api/server.go",
				EmbeddedStructs: []string{"BaseServer"},
				StructFields: []graph.StructField{
					{Name: "Port", Type: "int"},
				},
			},
			{Name: "BaseServer", Kind: graph.KindStruct, PackageName: "api", File: "pkg/api/server.go"},
			{Name: "Start", Kind: graph.KindMethod, Receiver: "*Server", PackageName: "api", Signature: "func (s *Server) Start() error", File: "pkg/api/server.go"},
			{Name: "stop", Kind: graph.KindMethod, Receiver: "*Server", PackageName: "api", Signature: "func (s *Server) stop()", File: "pkg/api/server.go"},
			{Name: "init", Kind: graph.KindFunction, PackageName: "api", Signature: "func init()", File: "pkg/api/server.go"},
		},
		Imports: []graph.ImportEdge{
			{FromFile: "pkg/api/server.go", FromPackage: "api", ImportPath: "net/http"},
		},
		Calls: []graph.CallEdge{
			{CallerName: "(*Server).Start", CalleeRaw: "http.ListenAndServe"},
			{CallerName: "main", CalleeRaw: "(*Server).Start"}, // 'Start' is called by 'main'
		},
		Mutations: []graph.MutationEdge{
			{Function: "(*Server).Start", Field: "Port", File: "pkg/api/server.go"},
		},
		Routes: []graph.HTTPRoute{
			{Method: "GET", Path: "/api/health", Handler: "healthHandler", File: "pkg/api/server.go"},
		},
		Errors: []graph.ErrorEdge{
			{Function: "(*Server).Start", Message: "failed to bind port", File: "pkg/api/server.go"},
		},
	}
}

func TestFocus(t *testing.T) {
	g := buildCoverageGraph()
	res := search.Focus(g, "api")
	if len(res) == 0 {
		t.Error("expected Focus to return results for 'api'")
	}
}

func TestSkeleton(t *testing.T) {
	g := buildCoverageGraph()
	out := search.Skeleton(g)
	if out == "" || !strings.Contains(out, "type Server struct") {
		t.Error("expected Skeleton to contain Server struct")
	}
}

func TestMutate(t *testing.T) {
	g := buildCoverageGraph()
	res := search.Mutate(g, "Port")
	if len(res) == 0 {
		t.Error("expected Mutate to return results for 'Port'")
	}
}

func TestImpact(t *testing.T) {
	g := buildCoverageGraph()
	res := search.Impact(g, "(*Server).Start")
	if len(res) == 0 {
		t.Error("expected Impact to return results for '(*Server).Start'")
	}
}

func TestTrace(t *testing.T) {
	g := buildCoverageGraph()
	res := search.Trace(g, "bind port")
	if len(res) == 0 {
		t.Error("expected Trace to return results for 'bind port'")
	}
}

func TestRoutes(t *testing.T) {
	g := buildCoverageGraph()
	res := search.Routes(g)
	if len(res) == 0 {
		t.Error("expected Routes to return results")
	}
}

func TestFields(t *testing.T) {
	g := buildCoverageGraph()
	res := search.Fields(g, "Server")
	if len(res) == 0 {
		t.Error("expected Fields to return results for 'Server'")
	}
}

func TestOrphans(t *testing.T) {
	g := buildCoverageGraph()
	res := search.Orphans(g)
	// stop is an orphan because it has no incoming calls. (Start is called by main)
	if len(res) == 0 {
		t.Error("expected Orphans to return results")
	}
}

func TestErrors(t *testing.T) {
	g := buildCoverageGraph()
	res := search.Errors(g, "")
	if len(res) == 0 {
		t.Error("expected Errors to return results")
	}
}

func TestPublic(t *testing.T) {
	g := buildCoverageGraph()
	res := search.Public(g, "api")
	if len(res) == 0 {
		t.Error("expected Public to return results for 'api'")
	}
}

func TestEmbeds(t *testing.T) {
	g := buildCoverageGraph()
	res := search.Embeds(g, "BaseServer")
	if len(res) == 0 {
		t.Error("expected Embeds to return results for 'BaseServer'")
	}
}
