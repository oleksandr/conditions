conditions
==========

This package offers a parser of a simple conditions specification language 
(reduced set of arithmetic/logical operations). It was created for Flow-Based Programming components that 
require configuration to perform some operations on the data received from multiple input ports.

Example:
```
import "github.com/oleksandr/conditions"


// Our condition to check
s := "($0 > 0.45) AND ($1 == \"ON\" OR $2 == \"ACTIVE\") AND $3 == false"

// Parse the condition language and get expression
p := NewParser(strings.NewReader(s))
expr, err := p.Parse()
if err != nil {
   // ...
}

// Evaluate expression passing data for $vars
r, err := Evaluate(expr, 0.12, "OFF", "ACTIVE", true)
if err != nil {
   // ...
}

// r is false

```
