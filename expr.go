package main

import (
	"github.com/google/cel-go/cel"
)

type Type = cel.Type
type Option = cel.EnvOption

type Expr struct {
	program cel.Program
}

// ThisVariable 注册map类型的this变量，方便map入参的操作
func ThisVariable() Option {
	return cel.Variable("this", cel.MapType(cel.StringType, cel.AnyType))
}

func NewExpr(expression string, options ...Option) (*Expr, error) {
	env, err := cel.NewEnv(options...)
	if err != nil {
		return nil, err
	}

	ast, issues := env.Compile(expression)
	if issues.Err() != nil {
		return nil, issues.Err()
	}

	program, err := env.Program(ast)
	if err != nil {
		return nil, err
	}
	return &Expr{program: program}, nil
}

func (e *Expr) Eval(input any) (any, error) {
	result, _, err := e.program.Eval(input)
	if err != nil {
		return nil, err
	}
	return result.Value(), nil
}
