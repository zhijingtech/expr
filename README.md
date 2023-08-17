# expr

expr 是基于 github.com/google/cel-go 再次封装的表达式解析和执行的工具，简化了 cel-go 的用法。

已实现的特性：
- 表达式解析支持 [Common Expression Language (CEL)](https://github.com/google/cel-spec/blob/master/doc/intro.md)
- 表达式执行入参支持 Go 的基础类型或者 ProtoBuf 声明的类型

待实现的特性：
- 表达式解析支持自定义函数
- 表达式执行出参支持自定义解析器

## 用法

简单用法举例：
```go
	expr, err := NewExpr("this.value > 60", ThisVariable())
	if err != nil {
		panic(err)
	}

	result, err := expr.Eval(map[string]any{"this": map[string]any{"value": 50}})

	if err != nil {
		panic(err)
	}

	fmt.Println("result:", result)
	// result: false
```

自定义函数用法举例：
```go
expr, err := NewExpr("shake_hands(i,you)",
		Variable("i", StringType),
		Variable("you", StringType),
		Function("shake_hands",
			Overload("shake_hands_string_string", []*Type{StringType, StringType}, StringType,
				BinaryBinding(func(arg1, arg2 Val) Val {
					return String(fmt.Sprintf("%v and %v are shaking hands.\n", arg1, arg2))
				}),
			),
		))
	if err != nil {
		panic(err)
	}

	result, err := expr.Eval(map[string]any{
		"i":   "CEL",
		"you": func() Val { return String("world") },
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("result:", result)
	// result: CEL and world are shaking hands.
```
## 欢迎贡献

项目刚拉起，欢迎向 https://github.com/zhijingtech/expr 提交问题或者PR。
