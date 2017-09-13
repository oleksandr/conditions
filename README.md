# conditions

This package offers a parser of a simple conditions specification language (reduced set of arithmetic/logical operations). The package is mainly created for Flow-Based Programming components that require configuration to perform some operations on the data received from multiple input ports. But it can be used whereever you need externally define some logical conditions on the internal variables.

Additional credits for this package go to [Handwritten Parsers & Lexers in Go](http://blog.gopheracademy.com/advent-2014/parsers-lexers/) by Ben Johnson on [Gopher Academy blog](http://blog.gopheracademy.com) and [InfluxML package from InfluxDB repository](https://github.com/influxdb/influxdb/tree/master/influxql).

## Usage example 
```
package main

import (
    "fmt"
    "strings"

    "github.com/oleksandr/conditions"
)

func main() {
    // Our condition to check
    s := "($0 > 0.45) AND ($1 == `ON` OR $2 == \"ACTIVE\") AND $3 == false"

    // Parse the condition language and get expression
    p := conditions.NewParser(strings.NewReader(s))
    expr, err := p.Parse()
    if err != nil {
        // ...
    }

    // Evaluate expression passing data for $vars
    data := map[string]interface{}{"$0": 0.12, "$1": "OFF", "$2": "ACTIVE", "$3": false}
    r, err := conditions.Evaluate(expr, data)
    if err != nil {
        // ...
    }

    // r is false
    fmt.Println("Evaluation result:", r)
}

```

## Where do we use it?

Here is a diagram for a sample FBP flow (created using [FlowMaker](https://github.com/cascades-fbp/flowmaker)). You can see how we configure the ContextA process with a condition via IIP packet.

![](https://raw.githubusercontent.com/oleksandr/conditions/master/Example.png)
