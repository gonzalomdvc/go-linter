// Go Linter 13: Redundant if statement that returns a boolean value

package checks

import (
	"fmt"
	"go/ast"
	"go/token"
	"slices"

	"github.com/gonzalomdvc/go-linter/model"
	"github.com/gonzalomdvc/go-linter/packages"
	"github.com/gonzalomdvc/go-linter/ui"
)

func GL13(fset *token.FileSet, file *ast.File, state *packages.State) []model.Finding {
	var findings []model.Finding
	ast.Inspect(file, func(n ast.Node) bool {
		var visitedIfStmts []token.Pos
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			var ifStmt *ast.IfStmt
			afterIfStmt := false
			if funcDecl.Body == nil {
				return true
			}
			for _, stmt := range funcDecl.Body.List {
				if ifs, ok := stmt.(*ast.IfStmt); ok {
					if slices.Contains(visitedIfStmts, ifs.Pos()) {
						return true
					}
					visitedIfStmts = append(visitedIfStmts, ifs.Pos())
					if !bodyHasPlainReturn(ifs.Body) {
						return true
					}
					for {
						elseIf, ok := ifs.Else.(*ast.IfStmt)
						if !ok {
							break
						}
						if elseIf.Body == nil {
							break
						}
						if slices.Contains(visitedIfStmts, elseIf.Pos()) {
							return true
						}
						visitedIfStmts = append(visitedIfStmts, elseIf.Pos())
						if !bodyHasPlainReturn(elseIf.Body) {
							return true
						}
						ifs = elseIf
					}
					afterIfStmt = true
					ifStmt = ifs
				} else {
					returnStmt, ok := stmt.(*ast.ReturnStmt)
					if afterIfStmt && ok {
						if isPlainReturnStmt(returnStmt) {
							msg, err := ui.PrintAt(fset.Position(ifStmt.Cond.Pos()))
							if err != nil {
								fmt.Printf("Error producing message for GL13")
							}
							findings = append(findings, model.Finding{
								Position: fset.Position(returnStmt.Pos()),
								Message:  fmt.Sprintf("Returning a plain boolean value on condition can be simplified, use return %s", msg),
							})
							return true
						}
					}
					afterIfStmt = false
					return true

				}
			}
		}

		return true
	})

	return findings
}

// Takes in a statement's body and checks if the return is a plain boolean
// like 'return true' or 'return false'
func bodyHasPlainReturn(node ast.Node) bool {
	body, ok := node.(*ast.BlockStmt)
	if !ok {
		fmt.Println("Is not body")
		return false
	}
	if body.List == nil {
		return false
	}
	if len(body.List) > 1 {
		return false
	}
	return isPlainReturnStmt(body.List[0])
}

func isPlainReturnStmt(n ast.Stmt) bool {
	isPlainReturn := true
	if returnStmt, ok := n.(*ast.ReturnStmt); ok {
		if returnStmt.Results == nil || len(returnStmt.Results) > 1 {
			isPlainReturn = false
			return false
		}
		if ident, ok := returnStmt.Results[0].(*ast.Ident); ok {
			if !(ident.Name == "true") && !(ident.Name == "false") {
				isPlainReturn = false
			}
		} else {
			isPlainReturn = false
		}
	} else {
		isPlainReturn = false
	}
	return isPlainReturn
}
