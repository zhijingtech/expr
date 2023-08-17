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

func TestExampleExprWithFunc(t *testing.T) {
	expr, err := NewExpr("shake_hands(i,you)",
		Variable("i", StringType),
		Variable("you", StringType),
		Function("shake_hands",
			Overload("shake_hands_string_string", []*Type{StringType, StringType}, StringType,
				BinaryBinding(func(arg1, arg2 Val) Val {
					return String(fmt.Sprintf("%v and %v are shaking hands.\n", arg1, arg2))
				}),
			),
		))
	if err != nil {
		panic(err)
	}

	result, err := expr.Eval(map[string]any{
		"i":   "CEL",
		"you": func() Val { return String("world") },
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("result:", result)
	// result: CEL and world are shaking hands.
}
