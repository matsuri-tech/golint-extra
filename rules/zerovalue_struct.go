package rules

import (
	"fmt"
	"github.com/pkg/errors"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"strings"
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
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	recordFieldInspector := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	rc := RecordFields{}

	// レコードのフィールドを集める
	inspect.Preorder(recordFieldInspector, func(n ast.Node) {
		switch decl := n.(type) {
		case *ast.GenDecl:
			if decl.Tok.IsKeyword() && decl.Tok.String() == "type" {
				spec, ok := decl.Specs[0].(*ast.TypeSpec)
				if !ok {
					return
				}

				st, ok := spec.Type.(*ast.StructType)
				if !ok {
					return
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
	})

	commentMapFiles := map[string]ast.CommentMap{}
	for _, file := range pass.Files {
		commentMapFiles[pass.Fset.Position(file.Pos()).Filename] = ast.NewCommentMap(pass.Fset, file, file.Comments)
	}

	recordInitializerInspector := []ast.Node{
		(*ast.CompositeLit)(nil),
	}

	var result []error

	// 全てのexpressionについてチェック
	inspect.Preorder(recordInitializerInspector, func(n ast.Node) {
		for _, comment := range commentMapFiles[pass.Fset.Position(n.Pos()).Filename].Filter(n).Comments() {
			if strings.HasPrefix(comment.Text(), "@ignore-golint-extra") {
				return
			}
		}

		switch expr := n.(type) {
		case *ast.CompositeLit:
			p, ok := expr.Type.(*ast.Ident)
			if !ok {
				return
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

				return
			}

			// KeyValueExprが続くときだけ処理する
			if len(expr.Elts) == 0 {
				return
			}

			_, ok = expr.Elts[0].(*ast.KeyValueExpr)
			if !ok {
				return
			}

			var keys []string
			for _, e := range expr.Elts {
				keys = append(keys, e.(*ast.KeyValueExpr).Key.(*ast.Ident).Name)
			}

			expected := rc[p.String()]
			if df := diff(expected, keys); len(df) != 0 {
				result = append(result, errors.New(fmt.Sprintf("Incomplete struct found: %+v\nMissing fields: %+v\n", pass.Fset.Position(p.Pos()), df)))
				return
			}

			return
		}
	})

	if len(result) == 0 {
		return nil, nil
	}

	resultErrorString := ""
	for _, err := range result {
		resultErrorString = resultErrorString + err.Error()
	}

	return nil, errors.New(resultErrorString)
}

func NewZeroValueStruct() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "zerovalue_struct",
		Doc:  "zerovalue-struct",
		Run:  run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}
