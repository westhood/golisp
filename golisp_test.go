package golisp

import (
    "testing"
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

func TestVector(t *testing.T) {
    l := ArrayToList([]Expr{Atom("if"), IntExpr(1), IntExpr(2), IntExpr(3)})
    v := Vector([]Expr{IntExpr(1), IntExpr(2), l})
    assert(t, v.Eval(nil).String() == "[1 2 2]", "Vector Expr")
    assert(t, v[0] == IntExpr(1), "Vector index")
}

func TestLet(t *testing.T) {
    var l Expr
    l = ArrayToList([]Expr{Atom("let"), Vector([]Expr{Atom("a"), IntExpr(1)}), Atom("a")})
    assert(t, l.Eval(nil) == IntExpr(1), "Let eval")

    l = ArrayToList([]Expr{Atom("let"),
            Vector([]Expr{Atom("a"), IntExpr(1), Atom("b"), IntExpr(2)}),
            ArrayToList([]Expr{Atom("do"), Atom("a"), Atom("b")})})
    assert(t, l.Eval(nil) == IntExpr(2), "Let eval")
}
