package main

import (
	"reflect"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/zhijingtech/expr/testdata"
)

func TestExpr_Eval(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		options    []Option
		input      any
		want       any
		wantErr    bool
	}{
		{
			name:       "expr on map return true",
			expression: "this.value == 1",
			options:    []Option{ThisVariable()},
			input:      map[string]any{"this": map[string]any{"value": 1}},
			want:       true,
		},
		{
			name:       "expr on map return false",
			expression: "this.value == 1",
			options:    []Option{ThisVariable()},
			input:      map[string]any{"this": map[string]any{"value": 2}},
			want:       false,
		},
		{
			name:       "expr on struct return true",
			expression: "(this.P2.X-this.P1.X) > 1.0",
			options: []Option{
				cel.Types(&testdata.Rectangle{}),
				cel.Variable("this", cel.ObjectType("testdata.Rectangle")),
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
				cel.Types(&testdata.Rectangle{}),
				cel.Variable("this", cel.ObjectType("testdata.Rectangle")),
			},
			input: map[string]any{"this": &testdata.Rectangle{
				P1: &testdata.Point{X: 1, Y: 2},
				P2: &testdata.Point{X: 3, Y: 4},
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewExpr(tt.expression, tt.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expr.NewExpr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := e.Eval(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expr.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
			options:    []Option{ThisVariable()},
			wantErr:    "ERROR: <input>:1:1: undeclared reference to 'dummy' (in container '')\n | dummy == 1\n | ^",
		},
		{
			name:       "expr wrong",
			expression: "dummy === 1",
			options:    []Option{ThisVariable()},
			wantErr:    "ERROR: <input>:1:9: Syntax error: token recognition error at: '= '\n | dummy === 1\n | ........^",
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
			options:    []Option{ThisVariable()},
			input:      map[string]any{},
			wantErr:    "no such attribute(s): this",
		},
		{
			name:       "variable not exist",
			expression: "this.dummy == 1",
			options:    []Option{ThisVariable()},
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
