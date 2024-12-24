package expr

import (
	"errors"
	"reflect"

	"github.com/google/cel-go/cel"
)

type (
	Env cel.Env
	// 定义一个接口，使用类型集来限制为基础类型
	Expr struct {
		p cel.Program
	}
)

var (
	ErrEnvNil = errors.New("env is nil")

	DefaultEnv, _ = NewEnv(UseThisVariable())

	ToSliceAny     = reflect.TypeOf([]any{})
	ToSliceInt     = reflect.TypeOf([]int{})
	ToSliceString  = reflect.TypeOf([]string{})
	ToMapStringAny = reflect.TypeOf(map[string]any{})
)

// UseThisVariable 注册map类型的this变量，方便在表达式中操作this数据
func UseThisVariable() Option {
	return cel.Variable("this", cel.MapType(cel.StringType, cel.DynType))
}

func WrapThisVariable(this map[string]any) map[string]any {
	return map[string]any{"this": this}
}

func NewEnv(opts ...Option) (*Env, error) {
	env, err := cel.NewEnv(opts...)
	return (*Env)(env), err
}

func (e *Env) Extend(opts ...Option) (*Env, error) {
	celEnv := (*cel.Env)(e)
	newEnv, err := celEnv.Extend(opts...)
	if err != nil {
		return nil, err
	}
	return (*Env)(newEnv), nil
}

func NewExpr(expression string, env ...*Env) (*Expr, error) {
	var _env *Env
	if len(env) == 0 {
		_env = DefaultEnv
	} else {
		_env = env[0]
	}
	celEnv := (*cel.Env)(_env)
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

func (e *Expr) Eval(input any, convertType ...reflect.Type) (any, error) {
	ev, _, err := e.p.Eval(input)
	if ev == nil || err != nil {
		return nil, err
	}
	if len(convertType) == 0 {
		return ev.Value(), nil
	}
	return ev.ConvertToNative(convertType[0])
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
