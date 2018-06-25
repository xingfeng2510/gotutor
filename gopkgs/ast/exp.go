package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func validOp(op token.Token) error {
	if op != token.EQL && op != token.GTR && op != token.LSS && op != token.LAND && op != token.LOR {
		return fmt.Errorf("unexpected operator: %s", op)
	}
	return nil
}

func parse(expr ast.Expr) error {
	switch t := expr.(type) {
	case *ast.Ident:
		fmt.Println(t.Name)
	case *ast.BasicLit:
		fmt.Println(t.Kind, t.Value)
	case *ast.BinaryExpr:
		fmt.Println(t.Op)
		if err := validOp(t.Op); err != nil {
			return err
		}
		if err := parse(t.X); err != nil {
			return err
		}
		if err := parse(t.Y); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected expression node type: %T", t)
	}
	return nil
}

func CheckExpr(x string) error {
	expr, err := parser.ParseExpr(x)
	if err != nil {
		return fmt.Errorf("parse expr %s failed", x)
	}
	if err := parse(expr); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := CheckExpr(`f1 == "abc" && f2 > 123`); err != nil {
		fmt.Println("err", err)
		return
	}
}
