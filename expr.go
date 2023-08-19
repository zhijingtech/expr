package expr

import (
	"github.com/google/cel-go/cel"
)

type Env cel.Env
type Expr struct {
	p cel.Program
}

// UseThisVariable 注册map类型的this变量，方便在表达式中操作this数据
func UseThisVariable() Option {
	return cel.Variable("this", cel.MapType(cel.StringType, cel.AnyType))
}

func NewEnv(opts ...Option) (*Env, error) {
	env, err := cel.NewEnv(opts...)
	return (*Env)(env), err
}

func (e *Env) NewExpr(expression string) (*Expr, error) {
	celEnv := (*cel.Env)(e)
	ast, issues := celEnv.Compile(expression)
	if issues.Err() != nil {
		return nil, issues.Err()
	}

	program, err := celEnv.Program(ast)
	if err != nil {
		return nil, err
	}
	return &Expr{p: program}, nil
}

func NewExpr(expression string, options ...Option) (*Expr, error) {
	env, err := NewEnv(options...)
	if err != nil {
		return nil, err
	}

	return env.NewExpr(expression)
}

func (e *Expr) Eval(input any) (any, error) {
	result, _, err := e.p.Eval(input)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.Value(), nil
}

// func (e *Expr) ContextEval(ctx context.Context, input any) (any, error) {
// 	result, _, err := e.p.ContextEval(ctx, input)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if result == nil {
// 		return nil, nil
// 	}
// 	return result.Value(), nil
// }
