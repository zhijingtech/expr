# expr

expr 是基于 github.com/google/cel-go 再次封装的表达式解析和执行的工具，简化了 cel-go 的用法。

已实现的特性：
- 表达式解析支持 [Common Expression Language (CEL)](https://github.com/google/cel-spec/blob/master/doc/intro.md)
- 表达式执行入参支持 Go 的基础类型或者 ProtoBuf 声明的类型
- 表达式解析支持自定义函数

待实现的特性：
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

内置的函数（macro）如下:

**has**
- **用途**：用于测试某个字段是否存在，避免需要将字段指定为字符串。
   - **用法举例**：`has(m.f)`
     - 如果对象 `m` 具有字段 `f`，则返回 `true`，否则返回 `false`。

**all**
- **用途**：将 `range.all(var, predicate)` 转换为一个理解式，确保范围内的所有元素都满足谓词条件。
   - **用法举例**：`range.all(x, x > 0)`
     - 如果范围内的所有元素都大于 0，则返回 `true`。

**exists**
- **用途**：将 `range.exists(var, predicate)` 转换为一个理解式，确保范围内至少有一个元素满足谓词条件。
   - **用法举例**：`range.exists(x, x > 0)`
     - 如果范围内至少有一个元素大于 0，则返回 `true`。

**existsOne**
- **用途**：将 `range.existsOne(var, predicate)` 转换为一个表达式，确保范围内恰好有一个元素满足谓词条件。
   - **用法举例**：`range.existsOne(x, x == 1)`
     - 如果范围内恰好有一个元素等于 1，则返回 `true`。

**map**
- **用途一**：将 `range.map(var, function)` 转换为一个理解式，对范围内的每个元素应用函数，生成一个新列表。
   - **用法举例**：`range.map(x, x * 2)`
     - 将范围内的每个元素乘以 2，并返回新列表。
- **用途二**：将 `range.map(var, predicate, function)` 转换为一个理解式，首先通过谓词过滤范围内的元素，然后应用转换函数生成一个新列表。
   - **用法举例**：`range.map(x, x > 0, x * 2)`
     - 过滤出大于 0 的元素，然后将这些元素乘以 2，并返回新列表。

**filter**
- **用途**：将 `range.filter(var, predicate)` 转换为一个理解式，过滤范围内的元素，生成一个满足谓词条件的新列表。
   - **用法举例**：`range.filter(x, x > 0)`
     - 返回一个新列表，其中包含范围内所有大于 0 的元素。

## 欢迎贡献

项目刚拉起，欢迎向 https://github.com/zhijingtech/expr 提交问题或者PR。
