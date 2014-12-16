package conditions

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	falseExpr = &BooleanLiteral{Val: false}
)

// Evaluate takes an expr and evaluates it using given args
func Evaluate(expr Expr, args ...interface{}) (bool, error) {
	if expr == nil {
		return false, fmt.Errorf("Provided expression is nil")
	}

	result, err := evaluateSubtree(expr, args...)
	if err != nil {
		return false, err
	}
	switch n := result.(type) {
	case *BooleanLiteral:
		return n.Val, nil
	}
	return false, fmt.Errorf("Unexpected result of the root expression: %#v", result)
}

// evaluateSubtree performs given expr evaluation recursively
func evaluateSubtree(expr Expr, args ...interface{}) (Expr, error) {
	if expr == nil {
		return falseExpr, fmt.Errorf("Provided expression is nil")
	}

	var (
		err    error
		lv, rv Expr
	)

	switch n := expr.(type) {
	case *ParenExpr:
		return evaluateSubtree(n.Expr, args...)
	case *BinaryExpr:
		lv, err = evaluateSubtree(n.LHS, args...)
		if err != nil {
			return falseExpr, err
		}
		rv, err = evaluateSubtree(n.RHS, args...)
		if err != nil {
			return falseExpr, err
		}
		return applyOperator(n.Op, lv, rv)
	case *VarRef:
		index, err := strconv.Atoi(strings.Replace(n.Val, "$", "", -1))
		if err != nil {
			return falseExpr, fmt.Errorf("Failed to resolve argument index: %s", err.Error())
		}
		if index >= len(args) {
			return falseExpr, fmt.Errorf("Not enough arguments provided. Number of arguments: %v. Requested element: %s", len(args), n.Val)
		}
		kind := reflect.TypeOf(args[index]).Kind()
		switch kind {
		case reflect.Int:
			return &NumberLiteral{Val: float64(args[index].(int))}, nil
		case reflect.Int32:
			return &NumberLiteral{Val: float64(args[index].(int32))}, nil
		case reflect.Int64:
			return &NumberLiteral{Val: float64(args[index].(int64))}, nil
		case reflect.Float32:
			return &NumberLiteral{Val: float64(args[index].(float32))}, nil
		case reflect.Float64:
			return &NumberLiteral{Val: float64(args[index].(float64))}, nil
		case reflect.String:
			return &StringLiteral{Val: args[index].(string)}, nil
		case reflect.Bool:
			return &BooleanLiteral{Val: args[index].(bool)}, nil
		}
		return falseExpr, fmt.Errorf("Unsupported argument %s type: %s", n.Val, kind)
	}

	return expr, nil
}

// applyOperator is a dispatcher of the evaluation according to operator
func applyOperator(op Token, l, r Expr) (*BooleanLiteral, error) {
	switch op {
	case AND:
		return applyAND(l, r)
	case OR:
		return applyOR(l, r)
	case EQ:
		return applyEQ(l, r)
	case NEQ:
		return applyNQ(l, r)
	case GT:
		return applyGT(l, r)
	case GTE:
		return applyGTE(l, r)
	case LT:
		return applyLT(l, r)
	case LTE:
		return applyLTE(l, r)
	}
	return &BooleanLiteral{Val: false}, fmt.Errorf("Unsupported operator: %s", op)
}

// applyAND applies && operation to l/r operands
func applyAND(l, r Expr) (*BooleanLiteral, error) {
	var (
		a, b bool
		err  error
	)
	a, err = getBoolean(l)
	if err != nil {
		return nil, err
	}
	b, err = getBoolean(r)
	if err != nil {
		return nil, err
	}
	return &BooleanLiteral{Val: (a && b)}, nil
}

// applyOR applies || operation to l/r operands
func applyOR(l, r Expr) (*BooleanLiteral, error) {
	var (
		a, b bool
		err  error
	)
	a, err = getBoolean(l)
	if err != nil {
		return nil, err
	}
	b, err = getBoolean(r)
	if err != nil {
		return nil, err
	}
	return &BooleanLiteral{Val: (a || b)}, nil
}

// applyEQ applies == operation to l/r operands
func applyEQ(l, r Expr) (*BooleanLiteral, error) {
	var (
		as, bs string
		an, bn float64
		ab, bb bool
		err    error
	)
	as, err = getString(l)
	if err == nil {
		bs, err = getString(r)
		if err != nil {
			return falseExpr, fmt.Errorf("Cannot compare string with non-string")
		}
		return &BooleanLiteral{Val: (as == bs)}, nil
	}
	an, err = getNumber(l)
	if err == nil {
		bn, err = getNumber(r)
		if err != nil {
			return falseExpr, fmt.Errorf("Cannot compare number with non-number")
		}
		return &BooleanLiteral{Val: (an == bn)}, nil
	}
	ab, err = getBoolean(l)
	if err == nil {
		bb, err = getBoolean(r)
		if err != nil {
			return falseExpr, fmt.Errorf("Cannot compare boolean with non-boolean")
		}
		return &BooleanLiteral{Val: (ab == bb)}, nil
	}
	return falseExpr, nil
}

// applyNQ applies != operation to l/r operands
func applyNQ(l, r Expr) (*BooleanLiteral, error) {
	var (
		as, bs string
		an, bn float64
		ab, bb bool
		err    error
	)
	as, err = getString(l)
	if err == nil {
		bs, err = getString(r)
		if err != nil {
			return falseExpr, fmt.Errorf("Cannot compare string with non-string")
		}
		return &BooleanLiteral{Val: (as != bs)}, nil
	}
	an, err = getNumber(l)
	if err == nil {
		bn, err = getNumber(r)
		if err != nil {
			return falseExpr, fmt.Errorf("Cannot compare number with non-number")
		}
		return &BooleanLiteral{Val: (an != bn)}, nil
	}
	ab, err = getBoolean(l)
	if err == nil {
		bb, err = getBoolean(r)
		if err != nil {
			return falseExpr, fmt.Errorf("Cannot compare boolean with non-boolean")
		}
		return &BooleanLiteral{Val: (ab != bb)}, nil
	}
	return falseExpr, nil
}

// applyGT applies > operation to l/r operands
func applyGT(l, r Expr) (*BooleanLiteral, error) {
	var (
		a, b float64
		err  error
	)
	a, err = getNumber(l)
	if err != nil {
		return nil, err
	}
	b, err = getNumber(r)
	if err != nil {
		return nil, err
	}
	return &BooleanLiteral{Val: (a > b)}, nil
}

// applyGTE applies >= operation to l/r operands
func applyGTE(l, r Expr) (*BooleanLiteral, error) {
	var (
		a, b float64
		err  error
	)
	a, err = getNumber(l)
	if err != nil {
		return nil, err
	}
	b, err = getNumber(r)
	if err != nil {
		return nil, err
	}
	return &BooleanLiteral{Val: (a >= b)}, nil
}

// applyLT applies < operation to l/r operands
func applyLT(l, r Expr) (*BooleanLiteral, error) {
	var (
		a, b float64
		err  error
	)
	a, err = getNumber(l)
	if err != nil {
		return nil, err
	}
	b, err = getNumber(r)
	if err != nil {
		return nil, err
	}
	return &BooleanLiteral{Val: (a < b)}, nil
}

// applyLTE applies <= operation to l/r operands
func applyLTE(l, r Expr) (*BooleanLiteral, error) {
	var (
		a, b float64
		err  error
	)
	a, err = getNumber(l)
	if err != nil {
		return falseExpr, err
	}
	b, err = getNumber(r)
	if err != nil {
		return falseExpr, err
	}
	return &BooleanLiteral{Val: (a <= b)}, nil
}

// getBoolean performs type assertion and returns boolean value or error
func getBoolean(e Expr) (bool, error) {
	switch n := e.(type) {
	case *BooleanLiteral:
		return n.Val, nil
	default:
		return false, fmt.Errorf("Literal is not a boolean: %v", n)
	}
}

// getString performs type assertion and returns string value or error
func getString(e Expr) (string, error) {
	switch n := e.(type) {
	case *StringLiteral:
		return n.Val, nil
	default:
		return "", fmt.Errorf("Literal is not a string: %v", n)
	}
}

// getNumber performs type assertion and returns float64 value or error
func getNumber(e Expr) (float64, error) {
	switch n := e.(type) {
	case *NumberLiteral:
		return n.Val, nil
	default:
		return 0, fmt.Errorf("Literal is not a number: %v", n)
	}
}
