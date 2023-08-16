package expr

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/zhijingtech/expr/testdata"
)

const (
	exprstr = "(current.P1.X-prev.P2.X) <= 1.0 || (current.P1.Y-prev.P2.Y) <= 1.0"
)

//go:generate protoc -I=. --go_out=. testdata/model.proto

func BenchmarkCelGoMap(b *testing.B) {
	expr, err := NewExpr(exprstr,
		cel.Variable("current", cel.MapType(cel.StringType, cel.AnyType)),
		cel.Variable("prev", cel.MapType(cel.StringType, cel.AnyType)),
	)
	assert.NoError(b, err)

	var prevMap = map[string]any{
		"P1": map[string]any{"X": 1, "Y": 2},
		"P2": map[string]any{"X": 3, "Y": 4},
	}
	var currentMap = map[string]any{
		"P1": map[string]any{"X": 5, "Y": 3},
		"P2": map[string]any{"X": 7, "Y": 5},
	}
	iparams := map[string]any{
		"current": currentMap,
		"prev":    prevMap,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := expr.Eval(iparams)
		assert.NoError(b, err)
		assert.True(b, result.(bool))
	}
}

func BenchmarkCelGoProtoBuf(b *testing.B) {
	expr, err := NewExpr(exprstr,
		cel.Types(&testdata.Rectangle{}),
		cel.Variable("prev", cel.ObjectType("testdata.Rectangle")),
		cel.Variable("current", cel.ObjectType("testdata.Rectangle")),
	)
	assert.NoError(b, err)

	var prev = &testdata.Rectangle{
		P1: &testdata.Point{X: 1, Y: 2},
		P2: &testdata.Point{X: 3, Y: 4},
	}
	var current = &testdata.Rectangle{
		P1: &testdata.Point{X: 5, Y: 3},
		P2: &testdata.Point{X: 7, Y: 5},
	}
	iparams := map[string]any{
		"prev":    prev,
		"current": current,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := expr.Eval(iparams)
		assert.NoError(b, err)
		assert.True(b, result.(bool))
	}
}
