package internal

import (
	"errors"
	"fmt"
	"go/ast"
	"log"
	"reflect"
)

func print(msg string, n ast.Node) error {
	buf := newBuffer()
	err := printNode(buf, n)
	if err != nil {
		return err
	}
	log.Println(msg, string(buf.buf.Bytes()))

	return err
}

func printNode(buff *buffer, n ast.Node) error {
	switch t := n.(type) {
	case *ast.Ident:
		buff.ws(t.Name)
	case *ast.AssignStmt:
		for i := 0; i < len(t.Lhs); i++ {
			lhs := t.Lhs[i]
			printNode(buff, lhs)
			if i+1 < len(t.Lhs) {
				buff.ws(", ")
			}
		}
		buff.ws(" ")
		buff.ws(t.Tok.String())
		buff.ws(" ")
		for i := 0; i < len(t.Rhs); i++ {
			rhs := t.Rhs[i]
			printNode(buff, rhs)
			if i+1 < len(t.Rhs) {
				buff.ws(", ")
			}
		}
	case *ast.CallExpr:
		printNode(buff, t.Fun)
		buff.ws("(")
		for i := 0; i < len(t.Args); i++ {
			arg := t.Args[i]
			printNode(buff, arg)
			if i+1 < len(t.Args) {
				buff.ws(", ")
			}
		}
		buff.ws(")")
	case *ast.ExprStmt:
		printNode(buff, t.X)
	case *ast.SelectorExpr:
		printNode(buff, t.X)
		buff.ws(".")
		printNode(buff, t.Sel)
	case *ast.BasicLit:
		buff.ws(t.Value)
	case *ast.ReturnStmt:
		buff.ws("return ")
		for i := 0; i < len(t.Results); i++ {
			r := t.Results[i]
			printNode(buff, r)
			if i+1 < len(t.Results) {
				buff.ws(", ")
			}
		}
	default:
		return errors.New(fmt.Sprintf("Unhandled node type %v", reflect.TypeOf(t)))
	}

	return nil
}
