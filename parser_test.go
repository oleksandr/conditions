package conditions

import (
	"strings"
	"testing"
)

var testData = []struct {
	cond   string
	args   []interface{}
	result bool
	isErr  bool
}{
	{"true", nil, true, false},
	{"false", nil, false, false},
	{"$0 > true", nil, false, true},
	{"$0 > true", []interface{}{43}, false, true},
	{"$0 > true", []interface{}{false}, false, true},
	{"$0 AND $1", []interface{}{true, true}, true, false},
	{"$0 AND $1", []interface{}{true, false}, false, false},
	{"$0 AND $1", []interface{}{false, true}, false, false},
	{"$0 AND $1", []interface{}{false, false}, false, false},
	{"$0 AND false", []interface{}{true}, false, false},
	{"56.43", nil, false, true},
	{"$5", nil, false, true},
	{"$0", []interface{}{true}, true, false},
	{"$0", []interface{}{false}, false, false},
	{"\"OFF\"", nil, false, true},
	{"$0 == \"OFF\"", []interface{}{"OFF"}, true, false},
	{"$0 > 10 AND $1 == \"OFF\"", []interface{}{14, "OFF"}, true, false},
	{"($0 > 10) AND ($1 == \"OFF\")", []interface{}{14, "OFF"}, true, false},
}

func TestInvalid(t *testing.T) {

}

func TestValid(t *testing.T) {
	var (
		expr Expr
		err  error
		r    bool
	)

	for _, td := range testData {
		t.Log("--------")
		t.Logf("Parsing: %s", td.cond)

		p := NewParser(strings.NewReader(td.cond))
		expr, err = p.Parse()
		t.Logf("Expression: %s", expr)
		if err != nil {
			t.Errorf("Unexpected error parsing expression: %s", td.cond)
			t.Error(err.Error())
			break
		}

		t.Log("Evaluating with:", td.args)
		r, err = Evaluate(expr, td.args...)
		if err != nil {
			if td.isErr {
				continue
			}
			t.Errorf("Unexpected error evaluating: %s", expr)
			t.Error(err.Error())
			break
		}
		if r != td.result {
			t.Errorf("Expected %v, received: %v", td.result, r)
		}
	}

	// Valid
	//s := "$1 > 3 OR (\"OFF\" == $0)"
	// s:= "true"
	// s:= "false"
	// s:= "$1 > true"
	// s:= "3 == true"
	// s:= "$1 > $2"
	//s := "$0 > 3 AND (78 > $0) AND ($0 >= -3 AND $1 < 20.3) OR ($2 > 10) AND ($1 != 44) AND ($2 <= 900) AND ($3 == \"ACTIVE\" OR $3 == \"IDLE\") OR $3 == $1 OR $3 == false"
	//s := "(P0 == -3) AND -100 >= P1"
	//s := "78 > P0 AND (P0 >= -3 AND P1 < 20.3) OR (P2 > 10) AND (P1 != 44) AND (P3 <= 900) AND (P3 == \"ACTIVE\" OR P3 == \"IDLE\") OR P3 == P1 OR P3 == false"

	// Invalid
	//s := "($1 >= -3 AND $1 < 20.3) OR ($2 >= 10s) AND ($3 == \"ACTIVE\" OR $3 == \"IDLE\") OR $3 == $1 OR $3 == false OR $5"

	/*
		p := NewParser(strings.NewReader(s))
		t.Log("Parsing...")
		t.Log(s)
		expr, err := p.Parse()
		if err != nil {
			t.Error(err.Error())
			t.FailNow()
		}

		t.Log("Evaluating...")
		t.Log(expr)
		//TODO: test case with empty args: matched, err := Evaluate(expr)
		matched, err := Evaluate(expr, "OFF", 56)
		t.Log("Analyzing...")
		if err != nil {
			t.Error(err.Error())
			t.FailNow()
		}
		if !matched {
			t.Errorf("Expected matched=true, but got %v", matched)
			t.Fail()
		}
	*/
}
