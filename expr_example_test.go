package expr

import (
	"fmt"
	"testing"
)

func TestExampleExprWithMap(t *testing.T) {
	expr, err := NewExpr("this.value > 60", UseThisVariable())
	if err != nil {
		panic(err)
	}

	result, err := expr.Eval(map[string]any{"this": map[string]any{"value": 50}})

	if err != nil {
		panic(err)
	}

	fmt.Println("result:", result)
	// result: false
}
