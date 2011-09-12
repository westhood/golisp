package golisp

import (
    "fmt"
    "strings"
)

type Context interface {
    Bind(name string, val Expr) bool
    Get(name string) (Expr, bool)
    Set(name string, val Expr)
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

// simple Context 
type SimpleContext struct {
    bindings map[string]Expr
    pre Context
}

func (context *SimpleContext) Bind(name string, val Expr) bool {
    if _, ok := (*context).bindings[name]; ok {
        return false
    }
    (*context).bindings[name] = val
    return true
}

func (context *SimpleContext) Get(name string) (val Expr , ok bool) {
    if val, ok = (*context).bindings[name]; ok {
        return
    }

    if (*context).pre != nil {
        val, ok = (*context).pre.Get(name)
        return
    }

    val, ok = nil, false
    return
}

func (context *SimpleContext) Set(name string, val Expr) {
    (*context).bindings[name] = val
    return
}


func NewContext(context Context) Context {
    return &SimpleContext{make(map[string]Expr), context}
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
        sa = append(sa, (*l).value.String())
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
        list = list.Cdr()
        switch string(form) {
            case "do": {
                return list.DoEval(context)
            }
            case "if": {
                return list.IfEval(context)
            }
            case "let": {
                return list.LetEval(context)
            }
        }
    }

    panic("Eval")
    // can't be here
    return nil
}

func (list *List) DoEval(context Context) Expr {
    seq := list
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
    if list.Len() != 3 {
        panic(fmt.Sprintf(
                "If Syntax Error: length of If form is %d",
                list.Len()))
    }

    l := list
    first := l.Car() // pred
    l = l.Cdr()
    second := l.Car() // true branch
    l = l.Cdr()
    third := l.Car() // false branch

    if bPred, ok := (first.Eval(context)).(Boolean); ok {
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

// (let [bindings] body)
func (list *List) LetEval(context Context) Expr {
    if list.Len() < 2 {
        panic(fmt.Sprintf(
                "Let Syntax Error: length of Let form is %d",
                list.Len()))
    }
    l := list

    var (
        binding Vector
        ok bool
    )

    if binding, ok = l.Car().(Vector); !ok {
         panic(fmt.Sprintf("Let Syntax Error: let requires a vector for its binding"))
    }
    if binding.Len() % 2 != 0 {
         panic(fmt.Sprintf("Let Syntax Error: let requires an even number of forms in binding vector"))
    }
    body := l.Cdr()

    newContext := NewContext(context)
    n := binding.Len()
    for i := 0; i < n/2 ; i++ {
        if name, ok := binding[2*i].(Atom); ok {
            val := binding[2*i+1].Eval(context)
            newContext.Bind(name.String(), val)
        } else {
            panic(fmt.Sprintf("Let Syntax Error: unsupported binding form"))
        }
    }

    return body.DoEval(newContext)
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


// Vector
type Vector []Expr

func (v Vector) Eval(context Context) Expr {
    val := []Expr{}
    for _, i := range v {
        val = append(val, i.Eval(context))
    }
    return Vector(val)
}

func (v Vector) Len() int {
    return len(v)
}

func (v Vector) String() string {
    sa := []string{}
    for _, i := range v {
        sa = append(sa, i.String())
    }
    return fmt.Sprintf("[%s]", strings.Join(sa, " "))
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

func (atom Atom) Eval(context Context) Expr {
    val, _ := context.Get(string(atom))
    return val
}


