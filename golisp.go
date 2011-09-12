package golisp

import (
    "fmt"
    "strings"
)

type Context interface {
}

type Expr interface {
    Eval(context Context) Expr
    String() string
}

type Proc interface {
    Apply(args []Expr, context Context) Expr
}

type Boolean interface {
    Bool() bool
}
// LIST
type List struct {
    value Expr
    next *List
}

func (list *List) Cons(value Expr) *List{
    newList := &List{value, list}
    return newList
}

func (list *List) Car() Expr {
    return (*list).value
}

func (list *List) Cdr() *List {
    return (*list).next
}

func (list *List) String() string {
    sa := []string{}
    for l := list; l != nil; l = (*l).next {
        sa = append(sa, fmt.Sprint((*l).value))
    }
    return fmt.Sprintf("(%s)", strings.Join(sa, " "))
}

func (list *List) Len() int {
    var n = 0
    for l:= list; l !=nil; l = l.Cdr() {
        n += 1
    }
    return n
}

func (list *List) Eval(context Context) Expr {
    begin := (*list).value
    if form, ok := begin.(Atom); ok {
        switch string(form) {
            case "do": {
                return list.DoEval(context)
            }
            case "if": {
                return list.IfEval(context)
            }
        }
    }
    return nil
}

func (list *List) DoEval(context Context) Expr {
    seq := list.Cdr()
    if seq == nil {
        return nil
    }

    var ret, e Expr
    for seq != nil {
        e = seq.Car()
        ret = e.Eval(context)
        seq = seq.Cdr()
    }
    return ret
}

func (list *List) IfEval(context Context) Expr {
    if list.Len() != 4 {
        panic(fmt.Sprintf(
                "If Syntax Error: length of If form is %d",
                list.Len()))
    }

    l := list.Cdr()
    first := l.Car()
    l = l.Cdr()
    second := l.Car() // true branch
    l = l.Cdr()
    third := l.Car() // false branch

    pred := first.Eval(context)

    if bPred, ok := pred.(Boolean); ok {
        if bPred.Bool() {
            return second.Eval(context)
        } else {
            return third.Eval(context)
        }
    } else {
        panic("pred is not a Boolean")
    }

    // can't be here
    return nil
}

func ArrayToList(array []Expr) *List {
    n := len(array)
    if n == 0 {
        return nil
    }
    l := &List{array[n-1], nil}
    for i:= n - 2; i >= 0; i-- {
        l = l.Cons(array[i])
    }
    return l
}


// Int
type IntExpr int

func (i IntExpr) String() string {
    return fmt.Sprint(int(i))
}

func (i IntExpr) Eval(_ Context) Expr {
    return i
}

func (i IntExpr) Bool() bool {
    if int(i) == 0 {
        return false
    }
    return true
}


// String
type StringExpr string

func (s StringExpr) String() string {
    return string(s)
}

func (s StringExpr) Eval(_ Context) Expr {
    return s
}


// Atom 
type Atom string

func (atom Atom) String() string {
    return string(atom)
}

func (atom Atom) Eval(_ Context) Expr {
    return atom
}


