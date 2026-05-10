package search

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ozgurcd/gograph/internal/graph"
)

// Skeleton returns a pseudo-Go string representing the structural API of the repository
// with all function bodies stripped.
func Skeleton(g *graph.Graph) string {
	var sb strings.Builder

	// Group symbols by package
	pkgSymbols := make(map[string][]graph.SymbolNode)
	for _, sym := range g.Symbols {
		pkgSymbols[sym.PackageName] = append(pkgSymbols[sym.PackageName], sym)
	}

	// Sort packages
	var pkgs []string
	for pkg := range pkgSymbols {
		pkgs = append(pkgs, pkg)
	}
	sort.Strings(pkgs)

	for _, pkg := range pkgs {
		sb.WriteString(fmt.Sprintf("package %s\n\n", pkg))

		// Sort symbols by Name inside package
		syms := pkgSymbols[pkg]
		sort.Slice(syms, func(i, j int) bool {
			return syms[i].Name < syms[j].Name
		})

		for _, sym := range syms {
			switch sym.Kind {
			case graph.KindStruct:
				sb.WriteString(fmt.Sprintf("type %s struct {\n", sym.Name))
				for _, emb := range sym.EmbeddedStructs {
					sb.WriteString(fmt.Sprintf("\t%s\n", emb))
				}
				for _, field := range sym.StructFields {
					if field.Tag != "" {
						sb.WriteString(fmt.Sprintf("\t%s %s `%s`\n", field.Name, field.Type, field.Tag))
					} else {
						sb.WriteString(fmt.Sprintf("\t%s %s\n", field.Name, field.Type))
					}
				}
				sb.WriteString("}\n\n")
			case graph.KindInterface:
				sb.WriteString(fmt.Sprintf("type %s interface {\n", sym.Name))
				// We don't have ordered interface methods, but we can print them
				var methods []string
				for m, sig := range sym.InterfaceMethods {
					methods = append(methods, fmt.Sprintf("\t%s%s", m, strings.TrimPrefix(sig, m)))
				}
				sort.Strings(methods)
				for _, m := range methods {
					sb.WriteString(m + "\n")
				}
				sb.WriteString("}\n\n")
			case graph.KindFunction:
				sb.WriteString(fmt.Sprintf("%s\n", sym.Signature))
			case graph.KindMethod:
				sb.WriteString(fmt.Sprintf("%s\n", sym.Signature))
			}
		}
		sb.WriteString("\n// " + strings.Repeat("-", 40) + "\n\n")
	}

	return sb.String()
}
