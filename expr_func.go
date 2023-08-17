package expr

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
)

type Val = ref.Val
type UnaryOp = functions.UnaryOp
type BinaryOp = functions.BinaryOp
type FunctionOp = functions.FunctionOp

type Null = types.Null
type Bool = types.Bool
type Int = types.Int
type Uint = types.Uint
type Bytes = types.Bytes
type Double = types.Double
type String = types.String
type Duration = types.Duration
type Unknown = types.Unknown
type Error = types.Error

func UnaryBinding(binding UnaryOp) cel.OverloadOpt {
	return cel.UnaryBinding(binding)
}

func BinaryBinding(binding BinaryOp) cel.OverloadOpt {
	return cel.BinaryBinding(binding)
}

func FunctionBinding(binding FunctionOp) cel.OverloadOpt {
	return cel.FunctionBinding(binding)
}

func Overload(overloadID string, args []*Type, resultType *Type, opts ...cel.OverloadOpt) cel.FunctionOpt {
	return cel.Overload(overloadID, args, resultType, opts...)
}

func Function(name string, opts ...cel.FunctionOpt) Option {
	return cel.Function(name, opts...)
}
