package expr

import (
	"github.com/google/cel-go/cel"
)

type Expr struct {
	program cel.Program
}

// UseThisVariable 注册map类型的this变量，方便在表达式中操作this数据
func UseThisVariable() Option {
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
