# expr

expr 是基于 github.com/google/cel-go 再次封装的表达式解析和执行的工具，简化了 cel-go 的用法。

已实现的特性：
- 表达式解析支持 [Common Expression Language (CEL)](https://github.com/google/cel-spec/blob/master/doc/intro.md)
- 表达式执行入参支持 Go 的基础类型或者 ProtoBuf 声明的类型

待实现的特性：
- 表达式解析支持自定义函数
- 表达式执行出参支持自定义解析器
