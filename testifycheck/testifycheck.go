//go:build !solution

package testifycheck

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

var Analyzer = &analysis.Analyzer{
	Name: "require",
	Doc:  "Some doc here",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	names := map[string]string{
		"Nil":     "NoError",
		"Nilf":    "NoErrorf",
		"NotNil":  "Error",
		"NotNilf": "Errorf",
	}

	for i, file := range pass.Files {
		if i == 1 {
			continue
		}
		processFile(pass, file, names)
	}

	return nil, nil
}

func processFile(pass *analysis.Pass, file *ast.File, names map[string]string) {
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return true
		}

		if _, ok := n.(*ast.ReturnStmt); ok {
			return false
		}

		ce, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		processCallExpr(pass, ce, names)
		return true
	})
}

func processCallExpr(pass *analysis.Pass, ce *ast.CallExpr, names map[string]string) {
	fn, _ := typeutil.Callee(pass.TypesInfo, ce).(*types.Func)
	if fn == nil || !isRequireOrAssert(fn) {
		return
	}

	if !hasErrorArgument(pass, ce) {
		return
	}

	if replacementName, ok := names[fn.Name()]; ok {
		pass.Reportf(ce.Pos(), "use %s.%s instead of comparing error to nil", fn.Pkg().Name(), replacementName)
	}
}

func isRequireOrAssert(fn *types.Func) bool {
	pkgName := fn.Pkg().Name()
	return pkgName == "require" || pkgName == "assert"
}

func hasErrorArgument(pass *analysis.Pass, ce *ast.CallExpr) bool {
	isErr := func(expr ast.Expr) bool {
		t := pass.TypesInfo.TypeOf(expr)
		if t == nil {
			return false
		}

		intf, ok := t.Underlying().(*types.Interface)
		if !ok {
			return false
		}

		return intf.NumMethods() == 1 && intf.Method(0).FullName() == "(error).Error"
	}

	argsLen := len(ce.Args)
	if argsLen < 1 {
		return false
	}
	return isErr(ce.Args[0]) || (argsLen > 1 && isErr(ce.Args[1]))
}
