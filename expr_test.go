package expr

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zhijingtech/expr/testdata"
)

type mapped map[string]any

func TestExpr_Eval(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		options    []Option
		input      any
		want       any
	}{
		{
			name:       "expr on map return true",
			expression: "this.value == 1",
			options:    []Option{UseThisVariable()},
			input:      map[string]any{"this": map[string]any{"value": 1}},
			want:       true,
		},
		{
			name:       "expr on map return false",
			expression: "this.value == 1",
			options:    []Option{UseThisVariable()},
			input:      map[string]any{"this": map[string]any{"value": 2}},
			want:       false,
		},
		{
			name:       "expr on struct return true",
			expression: "(this.P2.X-this.P1.X) > 1.0",
			options: []Option{
				Types(&testdata.Rectangle{}),
				Variable("this", ObjectType("testdata.Rectangle")),
			},
			input: map[string]any{"this": &testdata.Rectangle{
				P1: &testdata.Point{X: 1, Y: 2},
				P2: &testdata.Point{X: 3, Y: 4},
			}},
			want: true,
		},
		{
			name:       "expr on struct return false",
			expression: "(this.P2.X-this.P1.X) < 1.0",
			options: []Option{
				Types(&testdata.Rectangle{}),
				Variable("this", ObjectType("testdata.Rectangle")),
			},
			input: map[string]any{"this": &testdata.Rectangle{
				P1: &testdata.Point{X: 1, Y: 2},
				P2: &testdata.Point{X: 3, Y: 4},
			}},
			want: false,
		},
		{
			name:       "expr on struct return false",
			expression: "(this.P2.X-this.P1.X) < 1.0",
			options:    []Option{UseThisVariable()},
			input: map[string]any{"this": mapped{
				"P1": mapped{"X": 1, "Y": 2},
				"P2": mapped{"X": 3, "Y": 4},
			}},
			want: false,
		},
		// 自定义函数
		{
			name:       "expr custom function",
			expression: "distance(this.X, this.Y) < 1.0",
			options: []Option{
				UseThisVariable(),
				Function("distance",
					Overload("distance_d_d_d", []*Type{DoubleType, DoubleType}, DoubleType, BinaryBinding(func(arg1, arg2 Val) Val {
						d1 := arg1.(Double)
						d2 := arg2.(Double)
						dis := d1 - d2
						if dis >= 0.0 {
							return Double(dis)
						} else {
							return Double(-dis)
						}
					})))},
			input: map[string]any{"this": map[string]any{"X": 3.0, "Y": 3.5}},
			want:  true,
		},
		{
			name:       "expr custom function -2",
			expression: "distance(this.X, this.Y) > 1.0",
			options: []Option{
				UseThisVariable(),
				Function("distance",
					Overload("distance_d_d_d", []*Type{DoubleType, DoubleType}, DoubleType, BinaryBinding(func(arg1, arg2 Val) Val {
						d1 := arg1.(Double)
						d2 := arg2.(Double)
						dis := d1 - d2
						if dis < 0.0 {
							dis = -dis
						}
						return Double(dis)
					})))},
			input: map[string]any{"this": map[string]any{"X": 3.0, "Y": 3.5}},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewExpr(tt.expression, tt.options...)
			assert.NoError(t, err)
			got, err := e.Eval(tt.input)
			assert.NoError(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expr.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExpr_NewExpr_Err(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		options    []Option
		wantErr    string
	}{
		{
			name:       "variable not exist",
			expression: "dummy == 1",
			options:    []Option{UseThisVariable()},
			wantErr:    "ERROR: <input>:1:1: undeclared reference to 'dummy' (in container '')\n | dummy == 1\n | ^",
		},
		{
			name:       "expr wrong",
			expression: "dummy === 1",
			options:    []Option{UseThisVariable()},
			wantErr:    "ERROR: <input>:1:9: Syntax error: token recognition error at: '= '\n | dummy === 1\n | ........^",
		},
		// 自定义函数
		{
			name:       "function not exist",
			expression: "dummy(this.A)",
			options:    []Option{UseThisVariable()},
			wantErr:    "ERROR: <input>:1:6: undeclared reference to 'dummy' (in container '')\n | dummy(this.A)\n | .....^",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewExpr(tt.expression, tt.options...)
			assert.Nil(t, e)
			assert.Equal(t, tt.wantErr, err.Error())
		})
	}
}

func TestExpr_Eval_Err(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		options    []Option
		input      any
		wantErr    string
	}{
		{
			name:       "this not exist",
			expression: "this.dummy == 1",
			options:    []Option{UseThisVariable()},
			input:      map[string]any{},
			wantErr:    "no such attribute(s): this",
		},
		{
			name:       "variable not exist",
			expression: "this.dummy == 1",
			options:    []Option{UseThisVariable()},
			input:      map[string]any{"this": map[string]any{"value": 1}},
			wantErr:    "no such key: dummy",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewExpr(tt.expression, tt.options...)
			assert.NoError(t, err)

			got, err := e.Eval(tt.input)
			assert.Nil(t, got)
			assert.Equal(t, tt.wantErr, err.Error())
		})
	}
}

func TestExpr_ContextEval(t *testing.T) {
	options := []Option{
		UseThisVariable(),
		Function("sleep",
			Overload("sleep_int_null", []*Type{IntType}, NullType, UnaryBinding(func(arg Val) Val {
				t := arg.(Int).Value().(int)
				time.Sleep(time.Millisecond * time.Duration(t))
				return nil
			})))}

	env, err := NewEnv(options...)
	assert.NoError(t, err)

	expr1, err := env.NewExpr("this.milliseconds")
	assert.NoError(t, err)
	assert.NotNil(t, expr1)
	got, err := expr1.Eval(map[string]any{"this": map[string]any{"milliseconds": 200}})
	assert.NoError(t, err)
	assert.Equal(t, int64(200), got)

	// expr2, err := env.NewExpr("sleep(this.milliseconds)")
	// assert.NoError(t, err)
	// assert.NotNil(t, expr2)
	// ctx, canel := context.WithTimeout(context.Background(), time.Microsecond*100)
	// defer canel()
	// got, err = expr2.ContextEval(ctx, map[string]any{"this": map[string]any{"milliseconds": 200}})
	// assert.NoError(t, err)
	// assert.Nil(t, got)
}
