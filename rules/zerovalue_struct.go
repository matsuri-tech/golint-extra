package rules

import (
	"fmt"
	"github.com/pkg/errors"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
)

func diff(xs []string, ys []string) []string {
	var rs []string
	for _, x := range xs {
		found := false

		for _, y := range ys {
			if x == y {
				found = true
				break
			}
		}

		if !found {
			rs = append(rs, x)
		}
	}

	return rs
}

type RecordFields map[string][]string

// 再帰呼び出ししてるのでStackOverflowするかも
func checkExpr(fset token.FileSet, rc RecordFields, expr ast.Expr) error {
	switch expr := expr.(type) {
	case *ast.CompositeLit:
		p, ok := expr.Type.(*ast.Ident)
		if !ok {
			return nil
		}
		if p == nil {
			/* 次のようなものが見られるがこれは何？
				0  *ast.CompositeLit {
				1  .  Type: *ast.SelectorExpr {
				2  .  .  X: *ast.Ident {
				3  .  .  .  NamePos: -
				4  .  .  .  Name: "testdeps"
				5  .  .  }
				6  .  .  Sel: *ast.Ident {
				7  .  .  .  NamePos: -
				8  .  .  .  Name: "TestDeps"
				9  .  .  }
			   10  .  }
			   11  .  Lbrace: -
			   12  .  Rbrace: -
			   13  .  Incomplete: false
			   14  }
			*/

			return nil
		}

		// KeyValueExprが続くときだけ処理する
		_, ok = expr.Elts[0].(*ast.KeyValueExpr)
		if !ok {
			return nil
		}

		var keys []string
		for _, e := range expr.Elts {
			keys = append(keys, e.(*ast.KeyValueExpr).Key.(*ast.Ident).Name)
		}

		expected := rc[p.String()]
		if df := diff(expected, keys); len(df) != 0 {
			return errors.New(fmt.Sprintf("Incomplete struct found: %+v\nMissing fields: %+v\n", fset.Position(p.Pos()), df))
		}

		return errors.New(fmt.Sprintf("%+v", keys))
	case *ast.ParenExpr:
		if err := checkExpr(fset, rc, expr.X); err != nil {
			return err
		}
	case *ast.CallExpr:
		for _, arg := range expr.Args {
			if err := checkExpr(fset, rc, arg); err != nil {
				return err
			}
		}
	}

	return nil
}

func run(pass *analysis.Pass) (interface{}, error) {
	rc := RecordFields{}

	// レコードのフィールドを集める
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *ast.GenDecl:
				if decl.Tok.IsKeyword() && decl.Tok.String() == "type" {
					spec, ok := decl.Specs[0].(*ast.TypeSpec)
					if !ok {
						continue
					}

					st, ok := spec.Type.(*ast.StructType)
					if !ok {
						continue
					}

					var fields []string
					for _, f := range st.Fields.List {
						// structでnameが複数になることあるのだろうか…
						for _, n := range f.Names {
							fields = append(fields, n.String())
						}
					}

					rc[spec.Name.String()] = fields
				}
			}
		}
	}

	// 全てのexpressionについてチェック
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *ast.FuncDecl:
				for _, stmt := range decl.Body.List {
					switch stmt := stmt.(type) {
					case *ast.ExprStmt:
						if err := checkExpr(*pass.Fset, rc, stmt.X); err != nil {
							return nil, err
						}
					case *ast.AssignStmt:
						for _, expr := range stmt.Rhs {
							if err := checkExpr(*pass.Fset, rc, expr); err != nil {
								return nil, err
							}
						}
					default:
						return nil, nil
					}
				}
			}
		}
	}

	return nil, nil
}

func NewZeroValueStruct() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "zerovalue_struct",
		Doc:  "zerovalue-struct",
		Run:  run,
	}
}
