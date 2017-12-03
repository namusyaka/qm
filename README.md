# qm

Quine-McCluskey algorithm implementation written in Go

## Usage

```go
package main

import (
	"github.com/namusyaka/qm"
	"fmt"
)

func main() {
	// Registers boolean function.
	q := qm.New([]string{"A", "B", "C", "D"})

	// Calculates minified boolean function expressed as a sum of products.
	// And returns its complexity and result.
	complex, set := q.Solve([]int{4, 8, 10, 11, 12, 15}, []int{9, 14})
	fmt.Printf("complexity: %d\n", complex)
	fmt.Printf("set: %v\n", set)

	// Get boolean function.
	b := q.GetBoolFunc(set)
	fmt.Printf("boolean func: '%s'\n", b)
}
```
