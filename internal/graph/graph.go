// Package graph defines the core data model for the gograph tool.
package graph

import "time"

// Version is the schema version written into graph.json.
const Version = "1"

// Graph is the top-level data structure written to .gograph/graph.json.
type Graph struct {
	Version       string           `json:"version"`
	GeneratedAt   time.Time        `json:"generated_at"`
	Root          string           `json:"root"`
	Packages      []PackageNode    `json:"packages"`
	Files         []FileNode       `json:"files"`
	Symbols       []SymbolNode     `json:"symbols"`
	Imports       []ImportEdge     `json:"imports"`
	Calls         []CallEdge       `json:"calls"`
	EnvReads      []EnvRead        `json:"env_reads"`
	Dependencies  []Dependency     `json:"dependencies"`
	Routes        []HTTPRoute      `json:"routes,omitempty"`
	ImportEdges   []ImportEdge     `json:"import_edges,omitempty"`
	SQLs          []SQLEdge        `json:"sqls,omitempty"`
	Errors        []ErrorEdge      `json:"errors,omitempty"`
	Concurrency   []ConcurrencyNode `json:"concurrency,omitempty"`
	TestEdges     []TestEdge       `json:"test_edges,omitempty"`
}

// SQLEdge represents an extracted SQL query.
type SQLEdge struct {
	Query    string `json:"query"`
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// ErrorEdge represents an extracted error message or panic.
type ErrorEdge struct {
	Message  string `json:"message"`
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// ConcurrencyNode represents a concurrency primitive found in code.
// Kind is one of: "goroutine", "channel_send", "channel_recv", "mutex_lock",
// "mutex_unlock", "rwmutex_lock", "rwmutex_unlock", "waitgroup_add",
// "waitgroup_wait", "once_do".
type ConcurrencyNode struct {
	Kind     string `json:"kind"`
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Detail   string `json:"detail,omitempty"`
}

// TestEdge links a *testing.T test function to the symbols it exercises.
type TestEdge struct {
	TestFunc string `json:"test_func"`
	Target   string `json:"target"`   // callee name that looks like a production symbol
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// Dependency represents a go.mod dependency.
type Dependency struct {
	Module  string `json:"module"`
	Version string `json:"version"`
}

// HTTPRoute represents an HTTP REST endpoint found in the AST.
type HTTPRoute struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Handler string `json:"handler"`
	File    string `json:"file"`
	Line    int    `json:"line"`
}

// PackageNode represents a Go package found in the repository.
type PackageNode struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	ImportPathBestEffort string   `json:"import_path_best_effort"`
	Dir                  string   `json:"dir"`
	Files                []string `json:"files"`
}

// FileNode represents a single .go source file.
type FileNode struct {
	ID          string `json:"id"`
	Path        string `json:"path"`
	PackageName string `json:"package_name"`
	Lines       int    `json:"lines"`
	Generated   bool   `json:"generated"`
}

// SymbolKind categorises a symbol.
type SymbolKind string

const (
	KindFunction  SymbolKind = "function"
	KindMethod    SymbolKind = "method"
	KindStruct    SymbolKind = "struct"
	KindInterface SymbolKind = "interface"
)

// SymbolNode represents a named symbol (function, method, struct, interface).
type SymbolNode struct {
	ID               string            `json:"id"`
	Kind             SymbolKind        `json:"kind"`
	Name             string            `json:"name"`
	Receiver         string            `json:"receiver,omitempty"`
	PackageName      string            `json:"package_name"`
	File             string            `json:"file"`
	Line             int               `json:"line"`
	EndLine          int               `json:"end_line"`
	Doc              string            `json:"doc,omitempty"`
	Signature        string            `json:"signature,omitempty"`
	MethodSignature  string            `json:"method_signature,omitempty"`
	InterfaceMethods map[string]string `json:"interface_methods,omitempty"`
	StructFields     []StructField     `json:"struct_fields,omitempty"`
	EmbeddedStructs  []string          `json:"embedded_structs,omitempty"`
}

// StructField represents a field inside a struct.
type StructField struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Tag  string `json:"tag,omitempty"`
}

// CallEdge records a call expression found inside a function/method body.
type CallEdge struct {
	CallerSymbolID string `json:"caller_symbol_id"`
	CallerName     string `json:"caller_name"`
	CalleeRaw      string `json:"callee_raw"`
	File           string `json:"file"`
	Line           int    `json:"line"`
}

// ImportEdge records an import statement in a file.
type ImportEdge struct {
	FromFile    string `json:"from_file"`
	FromPackage string `json:"from_package"`
	ImportPath  string `json:"import_path"`
	Alias       string `json:"alias,omitempty"`
}

// EnvRead records a detected environment variable read.
type EnvRead struct {
	Key      string `json:"key"`
	Accessor string `json:"accessor"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function,omitempty"`
}

