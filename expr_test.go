package expr

import (
	"math"
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
		wantErr    string
	}{
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
					}))),
			},
			input: map[string]any{"this": map[string]any{"X": 3.0, "Y": 3.5}},
			want:  true,
		},
		{
			name:       "expr custom function 2",
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
					}))),
			},
			input: map[string]any{"this": map[string]any{"X": 3.0, "Y": 3.5}},
			want:  false,
		},
		{
			name:       "expr custom function panic",
			expression: "distance(this.X, this.Y) > 1.0",
			options: []Option{
				UseThisVariable(),
				Function("distance",
					Overload("distance_d_d_d", []*Type{DoubleType, DoubleType}, DoubleType, BinaryBinding(func(arg1, arg2 Val) Val {
						panic("this is panic")
					}))),
			},
			input:   map[string]any{"this": map[string]any{"X": 3.0, "Y": 3.5}},
			want:    false,
			wantErr: "internal error: this is panic",
		},
		{
			name:       "expr custom function error",
			expression: "distance(this.X, this.Y) > 1.0",
			options: []Option{
				UseThisVariable(),
				Function("distance",
					Overload("distance_d_d_d", []*Type{DoubleType, DoubleType}, DoubleType, BinaryBinding(func(arg1, arg2 Val) Val {
						return NewErr("this is error")
					}))),
			},
			input:   map[string]any{"this": map[string]any{"X": 3.0, "Y": 3.5}},
			want:    false,
			wantErr: "this is error",
		},
		{
			name:       "expr custom function 3",
			expression: "this.P1.dis_x(this.P2) > 1.0",
			options: []Option{
				Types(&testdata.Rectangle{}),
				Variable("this", ObjectType("testdata.Rectangle")),
				Function("dis_x",
					MemberOverload("point_dis_point_double", []*Type{ObjectType("testdata.Point"), ObjectType("testdata.Point")}, DoubleType,
						BinaryBinding(func(lhs, rhs Val) Val {
							p1, _ := lhs.Value().(*testdata.Point)
							p2, _ := rhs.Value().(*testdata.Point)
							return Double(math.Abs(p1.X - p2.X))
						},
						)),
				),
			},
			input: map[string]any{
				"this": &testdata.Rectangle{
					P1: &testdata.Point{X: 1, Y: 2},
					P2: &testdata.Point{X: 3, Y: 4},
				},
			},
			want: true,
		},
		{
			name:       "expr on array",
			expression: "len(this.items) > 1",
			options: []Option{
				UseThisVariable(),
				Function("len",
					Overload("len_Point_Int", []*Type{ListType(AnyType)}, IntType,
						UnaryBinding(func(arg Val) Val {
							v, _ := arg.Value().([]*testdata.Point)
							return Int(len(v))
						}))),
			},
			input: map[string]any{"this": map[string]any{"items": []*testdata.Point{{X: 1, Y: 2}, {X: 3, Y: 4}}}},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := NewEnv(tt.options...)
			assert.NoError(t, err)
			ex, err := NewExpr(tt.expression, env)
			assert.NoError(t, err)
			got, err := ex.Eval(tt.input)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.EqualError(t, err, tt.wantErr)
				assert.Nil(t, got)
			}
		})
	}
}

func TestExpr_EvalMap(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		input      map[string]any
		want       any
		wantErr    string
	}{
		{
			name:       "expr on map check int return true",
			expression: "this.value == 1",
			input:      map[string]any{"value": 1},
			want:       true,
		},
		{
			name:       "expr on map check int return false",
			expression: "this.value == 1",
			input:      map[string]any{"value": 2},
			want:       false,
		},
		{
			name:       "expr on map check string return true",
			expression: "this.value !=\"\"",
			input:      map[string]any{"value": "1"},
			want:       true,
		},
		{
			name:       "expr on map check string return false 1",
			expression: "this.value !=\"\"",
			input:      map[string]any{"value": ""},
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex, err := NewExpr(tt.expression)
			assert.NoError(t, err)
			got, err := ex.Eval(WrapThisVariable(tt.input))
			if tt.wantErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.EqualError(t, err, tt.wantErr)
				assert.Nil(t, got)
			}
		})
	}
}

func TestExpr_EvalAny(t *testing.T) {
	e, err := NewExpr("this.value")
	assert.NoError(t, err)

	input := map[string]any{"this": map[string]any{"value": 1}}
	got, err := e.Eval(input)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), got)

	input = map[string]any{"this": map[string]any{"value": "1"}}
	got, err = e.Eval(input)
	assert.NoError(t, err)
	assert.Equal(t, "1", got)

	input = map[string]any{"this": map[string]any{"value": []int{1, 2}}}
	got, err = e.Eval(input)
	assert.NoError(t, err)
	assert.Equal(t, []int([]int{1, 2}), got)

	input = map[string]any{"this": map[string]any{"value": map[int]int{1: 1, 2: 2}}}
	got, err = e.Eval(input)
	assert.NoError(t, err)
	assert.Equal(t, map[int]int{1: 1, 2: 2}, got)

	//--

	e, err = NewExpr("2+this.value/100")
	assert.NoError(t, err)
	input = map[string]any{"this": map[string]any{"value": 1}}
	got, err = e.Eval(input)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), got)

	input = map[string]any{"this": map[string]any{"value": 301}}
	got, err = e.Eval(input)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), got)
}

func TestExpr_EvalCustomStruct(t *testing.T) {
	options := []Option{
		Types(&testdata.Rectangle{}),
		Variable("this", MapType(StringType, AnyType)), // ObjectType("testdata.Rectangle")
	}
	env, err := NewEnv(options...)
	assert.NoError(t, err)
	e, err := NewExpr("this.level.rect", env)
	assert.NoError(t, err)
	input := map[string]any{"this": map[string]any{
		"level": map[string]any{
			"rect": &testdata.Rectangle{P1: &testdata.Point{X: 1, Y: 2}},
		},
	}}
	got, err := e.Eval(input)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), got.(*testdata.Rectangle).P1.X)
}

func TestExpr_EvalDyn(t *testing.T) {
	e, err := NewExpr("1.0 < 2")
	assert.EqualError(t, err, "ERROR: <input>:1:5: found no matching overload for '_<_' applied to '(double, int)'\n | 1.0 < 2\n | ....^")
	assert.Nil(t, e)

	e, err = NewExpr("dyn(1.0) < 2")
	assert.NoError(t, err)
	input := map[string]any{}
	got, err := e.Eval(input)
	assert.NoError(t, err)
	assert.True(t, got.(bool))
}

func TestExpr_EvalMissingKey(t *testing.T) {
	e, err := NewExpr("this.v1 > 0")
	assert.NoError(t, err)
	input := map[string]any{}
	got, err := e.Eval(input)
	assert.EqualError(t, err, "no such attribute(s): this")
	assert.Nil(t, got)

	input = map[string]any{"this": map[string]any{}}
	got, err = e.Eval(input)
	assert.EqualError(t, err, "no such key: v1")
	assert.Nil(t, got)

	e, err = NewExpr("has(this.v1) && this.v1 > 0")
	assert.NoError(t, err)
	input = map[string]any{"this": map[string]any{}}
	got, err = e.Eval(input)
	assert.NoError(t, err)
	assert.Equal(t, false, got)

	input = map[string]any{"this": map[string]any{"v1": 1}}
	got, err = e.Eval(input)
	assert.NoError(t, err)
	assert.Equal(t, true, got)
}

func TestExpr_EvalJson(t *testing.T) {
	e, err := NewExpr(`{"a":this.a}`)
	assert.NoError(t, err)
	assert.NotNil(t, e)

	input := map[string]any{"this": map[string]any{"a": map[string]any{"b": 1}}}
	got, err := e.Eval(input)
	assert.NoError(t, err)
	assert.Equal(t, map[any]any(map[any]any{"a": map[string]any{"b": 1}}), got)
}

func TestExpr_EvalMicro(t *testing.T) {
	e, err := NewExpr(`this.array.filter(x, x > 1)`)
	assert.NoError(t, err)
	assert.NotNil(t, e)

	input := map[string]any{"this": map[string]any{"array": []int{1, 2, 3}}}
	got, err := e.Eval(input)
	assert.NoError(t, err)
	assert.Equal(t, []any{int64(2), int64(3)}, got)
}

func TestExpr_NewExpr_Err(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantErr    string
	}{
		{
			name:       "variable not exist",
			expression: "dummy == 1",
			wantErr:    "ERROR: <input>:1:1: undeclared reference to 'dummy' (in container '')\n | dummy == 1\n | ^",
		},
		{
			name:       "expr wrong",
			expression: "dummy === 1",
			wantErr:    "ERROR: <input>:1:9: Syntax error: token recognition error at: '= '\n | dummy === 1\n | ........^",
		},
		// 自定义函数
		{
			name:       "function not exist",
			expression: "dummy(this.A)",
			wantErr:    "ERROR: <input>:1:6: undeclared reference to 'dummy' (in container '')\n | dummy(this.A)\n | .....^",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewExpr(tt.expression)
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
			input:      map[string]any{},
			wantErr:    "no such attribute(s): this",
		},
		{
			name:       "variable not exist",
			expression: "this.dummy == 1",
			input:      map[string]any{"this": map[string]any{"value": 1}},
			wantErr:    "no such key: dummy",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewExpr(tt.expression)
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
			}))),
	}

	env, err := NewEnv(options...)
	assert.NoError(t, err)

	expr1, err := NewExpr("this.milliseconds", env)
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

func TestEnv_Extend(t *testing.T) {
	env, err := NewEnv()
	assert.NoError(t, err)
	assert.NotNil(t, env)

	// 未注册Variable和Function，基础计算
	expr, err := NewExpr("1 + 2", env)
	assert.NoError(t, err)
	assert.NotNil(t, expr)
	got, err := expr.Eval(map[string]any{})
	assert.NoError(t, err)
	assert.Equal(t, int64(3), got)

	// 未注册Variable和Function，调用未注册方法
	expr, err = NewExpr("ret(v)", env)
	assert.EqualError(t, err, "ERROR: <input>:1:4: undeclared reference to 'ret' (in container '')\n | ret(v)\n | ...^\nERROR: <input>:1:5: undeclared reference to 'v' (in container '')\n | ret(v)\n | ....^")
	assert.Nil(t, expr)

	// 未注册Variable和注册Function，调用注册方法
	env, err = env.Extend(Function("ret", Overload("ret_int_int", []*Type{IntType}, IntType, UnaryBinding(func(arg Val) Val {
		return arg
	}))))
	assert.NoError(t, err)
	assert.NotNil(t, env)
	expr, err = NewExpr("ret(v)", env)
	assert.EqualError(t, err, "ERROR: <input>:1:5: undeclared reference to 'v' (in container '')\n | ret(v)\n | ....^")
	assert.Nil(t, expr)

	// 注册Variable和注册Function，调用注册方法
	env, err = env.Extend(Variable("v", IntType))
	assert.NoError(t, err)
	assert.NotNil(t, env)
	expr, err = NewExpr("ret(v)", env)
	assert.NoError(t, err)
	assert.NotNil(t, expr)
	got, err = expr.Eval(map[string]any{"v": 3})
	assert.NoError(t, err)
	assert.Equal(t, int64(3), got)
}
