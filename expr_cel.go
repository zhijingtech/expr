package expr

import (
	"github.com/google/cel-go/cel"
)

var (
	AnyType       = cel.AnyType
	BoolType      = cel.BoolType
	BytesType     = cel.BytesType
	DoubleType    = cel.DoubleType
	DurationType  = cel.DurationType
	DynType       = cel.DynType
	IntType       = cel.IntType
	NullType      = cel.NullType
	StringType    = cel.StringType
	TimestampType = cel.TimestampType
	TypeType      = cel.TypeType
	UintType      = cel.UintType
)

type Type = cel.Type
type Option = cel.EnvOption

func Types(addTypes ...any) Option {
	return cel.Types(addTypes...)
}

func Variable(name string, t *Type) Option {
	return cel.Variable(name, t)
}

func ObjectType(typeName string) *Type {
	return cel.ObjectType(typeName)
}
