package conditions

import (
	"strings"
	"testing"
)

var invalidTestData = []string{
	"",
	"[] AND true",
	"A",
	"[var0] == DEMO",
	"[var0] == 'DEMO'",
	"![var0]",
	"[var0] <> `DEMO`",
}

var validTestData = []struct {
	cond   string
	args   map[string]interface{}
	result bool
	isErr  bool
}{
	{"true", nil, true, false},
	{"false", nil, false, false},
	{"false OR true OR false OR false OR true", nil, true, false},
	{"((false OR true) AND false) OR (false OR true)", nil, true, false},
	{"[var0]", map[string]interface{}{"var0": true}, true, false},
	{"[var0]", map[string]interface{}{"var0": false}, false, false},
	{"[var0] > true", nil, false, true},
	{"[var0] > true", map[string]interface{}{"var0": 43}, false, true},
	{"[var0] > true", map[string]interface{}{"var0": false}, false, true},
	{"[var0] and [var1]", map[string]interface{}{"var0": true, "var1": true}, true, false},
	{"[var0] AND [var1]", map[string]interface{}{"var0": true, "var1": false}, false, false},
	{"[var0] AND [var1]", map[string]interface{}{"var0": false, "var1": true}, false, false},
	{"[var0] AND [var1]", map[string]interface{}{"var0": false, "var1": false}, false, false},
	{"[var0] AND false", map[string]interface{}{"var0": true}, false, false},
	{"56.43", nil, false, true},
	{"[var5]", nil, false, true},
	{"[var0] > -100 AND [var0] < -50", map[string]interface{}{"var0": -75.4}, true, false},
	{"[var0]", map[string]interface{}{"var0": true}, true, false},
	{"[var0]", map[string]interface{}{"var0": false}, false, false},
	{"\"OFF\"", nil, false, true},
	{"`ON`", nil, false, true},
	{"[var0] == \"OFF\"", map[string]interface{}{"var0": "OFF"}, true, false},
	{"[var0] > 10 AND [var1] == \"OFF\"", map[string]interface{}{"var0": 14, "var1": "OFF"}, true, false},
	{"([var0] > 10) AND ([var1] == \"OFF\")", map[string]interface{}{"var0": 14, "var1": "OFF"}, true, false},
	{"([var0] > 10) AND ([var1] == \"OFF\") OR true", map[string]interface{}{"var0": 1, "var1": "ON"}, true, false},
	{"[foo][dfs] == true and [bar] == true", map[string]interface{}{"foo.dfs": true, "bar": true}, true, false},
	{"[foo][dfs][a] == true and [bar] == true", map[string]interface{}{"foo.dfs.a": true, "bar": true}, true, false},
	{"[@foo][a] == true and [bar] == true", map[string]interface{}{"@foo.a": true, "bar": true}, true, false},
	{"[foo][unknow] == true and [bar] == true", map[string]interface{}{"foo.dfs": true, "bar": true}, false, true},
	//XOR
	{"false XOR false", nil, false, false},
	{"false xor true", nil, true, false},
	{"true XOR false", nil, true, false},
	{"true xor true", nil, false, false},

	//NAND
	{"false NAND false", nil, true, false},
	{"false nand true", nil, true, false},
	{"true nand false", nil, true, false},
	{"true NAND true", nil, false, false},

	//IN
	{"[foo] in [foobar]", map[string]interface{}{"foo": "findme", "foobar": []string{"notme", "may", "findme", "lol"}}, true, false},
	//{`[foo] in ["hello", "world", "foo"]`, map[string]interface{}{"foo": "world"}, false, false},
	//{`[foo] in ["hello", "world", "foo"]`, map[string]interface{}{"foo": "monde"}, true, true},

	//NOT IN
	{"[foo] not in [foobar]", map[string]interface{}{"foo": "dontfindme", "foobar": []string{"notme", "may", "findme", "lol"}}, true, false},

	// =~
	{"[status] =~ /^5\\d\\d/", map[string]interface{}{"status": "500"}, true, false},
	{"[status] =~ /^4\\d\\d/", map[string]interface{}{"status": "500"}, false, false},

	// !~
	{"[status] !~ /^5\\d\\d/", map[string]interface{}{"status": "500"}, false, false},
	{"[status] !~ /^4\\d\\d/", map[string]interface{}{"status": "500"}, true, false},
}

func TestInvalid(t *testing.T) {

	var (
		expr Expr
		err  error
	)

	for _, cond := range invalidTestData {
		t.Log("--------")
		t.Logf("Parsing: %s", cond)

		p := NewParser(strings.NewReader(cond))
		expr, err = p.Parse()
		if err == nil {
			t.Error("Should receive error")
			break
		}
		if expr != nil {
			t.Error("Expression should nil")
			break
		}
	}
}

func TestValid(t *testing.T) {

	var (
		expr Expr
		err  error
		r    bool
	)

	for _, td := range validTestData {
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

		t.Log("Evaluating with: %#v", td.args)
		r, err = Evaluate(expr, td.args)
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
			break
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
