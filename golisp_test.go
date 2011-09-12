package golisp

import (
    "testing"
    // "fmt"
)

func assert(t *testing.T, pred bool, format string, args ...interface{}) {
    if !pred {
        t.Errorf(format, args...)
    }
}

func protect(t *testing.T, g func(), format string, args ...interface{}) {
	defer func() {
		if x := recover(); x != nil {
		} else {
            t.Errorf(format, args...)
        }
	}()
	g()
}

func TestList(t *testing.T) {
    l := ArrayToList([]Expr{IntExpr(1), IntExpr(2), IntExpr(3)})
    assert(t, l.String() == "(1 2 3)", "List to String")

    assert(t, l.Len() == 3, "List length")
    l = nil
    assert(t, l.Len() == 0, "List length")
}

func TestIfExpr(t *testing.T) {
    var l Expr
    l = ArrayToList([]Expr{Atom("if"), IntExpr(1), IntExpr(2), IntExpr(3)})
    assert(t, l.Eval(nil) == IntExpr(2), "If Expr")

    l = ArrayToList([]Expr{Atom("if"), IntExpr(0), IntExpr(2), IntExpr(3)})
    assert(t, l.Eval(nil) == IntExpr(3), "If Expr")

    l = ArrayToList([]Expr{Atom("if"), StringExpr("t"), IntExpr(2), IntExpr(3)})
    protect(t, func () { l.Eval(nil) }, "Pred is not a Boolean")
}
