package exitcheck

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const (
	targetPackage = "os"
	targetFunc    = "Exit"
)

var Analyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "Checks os.Exit from main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}

		for _, decl := range file.Decls {
			if fnDecl, ok := decl.(*ast.FuncDecl); ok && fnDecl.Name.Name == "main" {
				ast.Inspect(fnDecl.Body, func(node ast.Node) bool {
					if callExpr, callExprOk := node.(*ast.CallExpr); callExprOk {
						exitDetectInCallExpr(pass, callExpr)
					}
					return true
				})
			}
		}
	}

	return nil, nil
}

func exitDetectInCallExpr(pass *analysis.Pass, x *ast.CallExpr) {
	if s, ok := x.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := s.X.(*ast.Ident); ok && ident.Name == targetPackage && s.Sel.Name == targetFunc {
			pass.Reportf(ident.NamePos, "os.Exit call prohibited in main package, main function")
		}
	}
}
