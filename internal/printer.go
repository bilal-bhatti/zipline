package internal

import (
	"bytes"
	"go/ast"
	"log"
	"reflect"
)

func print(msg string, n ast.Node) {
	var buf bytes.Buffer
	printNode(&buf, n)
	log.Println(msg, string(buf.Bytes()))
}

func printNode(buff *bytes.Buffer, n ast.Node) {
	switch t := n.(type) {
	case *ast.Ident:
		buff.WriteString(t.Name)
	case *ast.AssignStmt:
		for i := 0; i < len(t.Lhs); i++ {
			lhs := t.Lhs[i]
			printNode(buff, lhs)
			if i+1 < len(t.Lhs) {
				buff.WriteString(", ")
			}
		}
		buff.WriteString(" ")
		buff.WriteString(t.Tok.String())
		buff.WriteString(" ")
		for i := 0; i < len(t.Rhs); i++ {
			rhs := t.Rhs[i]
			printNode(buff, rhs)
			if i+1 < len(t.Rhs) {
				buff.WriteString(", ")
			}
		}
	case *ast.CallExpr:
		printNode(buff, t.Fun)
		buff.WriteString("(")
		for i := 0; i < len(t.Args); i++ {
			arg := t.Args[i]
			printNode(buff, arg)
			if i+1 < len(t.Args) {
				buff.WriteString(", ")
			}
		}
		buff.WriteString(")")
	case *ast.ExprStmt:
		printNode(buff, t.X)
	case *ast.SelectorExpr:
		printNode(buff, t.X)
		buff.WriteString(".")
		printNode(buff, t.Sel)
	case *ast.BasicLit:
		buff.WriteString(t.Value)
	case *ast.ReturnStmt:
		buff.WriteString("return ")
		for i := 0; i < len(t.Results); i++ {
			r := t.Results[i]
			printNode(buff, r)
			if i+1 < len(t.Results) {
				buff.WriteString(", ")
			}
		}
	default:
		log.Println("Unhandled node type", reflect.TypeOf(t))
	}
}
