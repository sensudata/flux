// DO NOT EDIT: This file is autogenerated via the builtin command.

package kafka

import (
	flux "github.com/influxdata/flux"
	ast "github.com/influxdata/flux/ast"
)

func init() {
	flux.RegisterPackage(pkgAST)
}

var pkgAST = &ast.Package{
	BaseNode: ast.BaseNode{
		Errors: nil,
		Loc:    nil,
	},
	Files: []*ast.File{&ast.File{
		BaseNode: ast.BaseNode{
			Errors: nil,
			Loc: &ast.SourceLocation{
				End: ast.Position{
					Column: 11,
					Line:   3,
				},
				File:   "kafka.flux",
				Source: "package kafka\n\nbuiltin to",
				Start: ast.Position{
					Column: 1,
					Line:   1,
				},
			},
		},
		Body: []ast.Statement{&ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 11,
						Line:   3,
					},
					File:   "kafka.flux",
					Source: "builtin to",
					Start: ast.Position{
						Column: 1,
						Line:   3,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 11,
							Line:   3,
						},
						File:   "kafka.flux",
						Source: "to",
						Start: ast.Position{
							Column: 9,
							Line:   3,
						},
					},
				},
				Name: "to",
			},
		}},
		Imports:  nil,
		Metadata: "parser-type=go",
		Name:     "kafka.flux",
		Package: &ast.PackageClause{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   1,
					},
					File:   "kafka.flux",
					Source: "package kafka",
					Start: ast.Position{
						Column: 1,
						Line:   1,
					},
				},
			},
			Name: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   1,
						},
						File:   "kafka.flux",
						Source: "kafka",
						Start: ast.Position{
							Column: 9,
							Line:   1,
						},
					},
				},
				Name: "kafka",
			},
		},
	}},
	Package: "kafka",
	Path:    "kafka",
}
