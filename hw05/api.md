# MyCobra 包

`import "github.com/KianKw/myCobra"`

## 概述

MyCobra 是对 [Cobra](https://github.com/spf13/cobra) 的简单模仿，主要用于学习 Go 语言；

Cobra 是一个接口，用于创建 CLI 接口，类似于 git 、go 工具；

Cobra 是一个程序，可以快速开发基于 cobra 的其他程序。

## 目录 

* 变量

* 函数
    *  [func (cmd *Command) AddCommand(subcmd *Command)](#Command.AddCommand)
*  [func (cmd *Command) Execute() error](#Command.Execute)
    *  [func (cmd *Command, []string, error) Find(args []string) ( *Command, []string, error)](#Command.Find)
    *  [func (cmd *Command, []string) HelpFunc() func( *Command, []string)](#Command.HelpFunc)
    *  [func (cmd *Command) InitDefaultHelpCmd()](#Command.InitDefaultHelpCmd)
    *  [func (cmd *Command) Name() string](#Command.Name)
    *  [func (c *Command) PersistentFlags() *flag.FlagSet](#Command.PersistentFlags)
    *  [func (c *Command) Print_help()](#Command.Print_help)
    *  [func (cmd *Command) RemoveCommand()](#Command.RemoveCommand)
    *  [func (cmd *Command) Root() *Command](#Command.Root)
    *  [func (cmd *Command) Runnable() bool](#Command.Runnable)
    *  [func (cmd *Command) SetHelpFunc(f func( *Command, []string))](#Command.SetHelpFunc)

### 包内文件

[command.go](src/github.com/KianKw/command.go)

## 变量

type [Command struct](/src/github.com/KianKw/myCobra/command.go?s=113:424#L2) [¶](#Command)

* Use [string](/pkg/builtin/#string): 描述使用方法的单行信息.
* Short [string](/pkg/builtin/#string): Short 在 `help` 中输出的短描述.
* Long [string](/pkg/builtin/#string) : 在 `help` 中输出的长描述.
* Run func(cmd \*[Command](#Command), args \[\][string](/pkg/builtin/#string)) : 运行命令, 大多命令只需要实现这个函数.
* RunE func(cmd \*[Command](#Command), args \[\][string](/pkg/builtin/#string)) [error](/pkg/builtin/#error) : 运行命令并返回错误.
* args []string: args 通过 flags 解析得到的参数.
* parent *Command: parent 该命令的父命令.
* commands []*Command: commands 该命令支持的子命令列表.


## 方法

### func [ParseArgs](/src/github.com/KianKw/myCobra/command.go?s=1359:1400#L72) [¶](#ParseArgs)

```go
func ParseArgs(c *Command, args []string)
```

解析输入的命令的参数

### func (\*Command) [AddCommand](/src/github.com/KianKw/myCobra/command.go?s=2680:2731#L130) [¶](#Command.AddCommand)

```go
func (cmd *Command) AddCommand(subCmds ...*Command)
```

为该命令增加子命令.

### func (\*Command) [Execute](/src/github.com/KianKw/myCobra/command.go?s=428:463#L26) [¶](#Command.Execute)

```go
func (cmd *Command) Execute() error
```

调用 `Find` 函数找到要执行的目标命令，再调用 `execute` 执行该命令

### func (\*Command) [Find](/src/github.com/KianKw/myCobra/command.go?s=1694:1761#L90) [¶](#Command.Find)

```go
func (cmd *Command) Find(args []string) (*Command, []string, error)
```

从参数中找到要执行的子命令, 如果没有子命令则返回这个命令本身，如果找不到则返回错误

### func (\*Command) [HelpFunc](/src/github.com/KianKw/myCobra/command.go?s=4533:4588#L204) [¶](#Command.HelpFunc)

```go
func (cmd *Command) HelpFunc() func(*Command, []string)
```

HelpFunc 返回自身或父函数或默认 helpFunc.

### func (\*Command) [InitDefaultHelpCmd](/src/github.com/KianKw/myCobra/command.go?s=4919:4959#L223) [¶](#Command.InitDefaultHelpCmd)

```go
func (cmd *Command) InitDefaultHelpCmd()
```

初始化默认的 Help 命令

### func (\*Command) [Name](/src/github.com/KianKw/myCobra/command.go?s=4219:4252#L187) [¶](#Command.Name)

```go
func (cmd *Command) Name() string
```

返回 Command.Use 字段首单词作为该命令的名字.

### func (\*Command) [PersistentFlags](/src/github.com/KianKw/myCobra/command.go?s=4054:4103#L180) [¶](#Command.PersistentFlags)

```go
func (c *Command) PersistentFlags() *flag.FlagSet
```

返回全局的持久 Flag

### func (\*Command) [Print\_help](/src/github.com/KianKw/myCobra/command.go?s=3255:3285#L155) [¶](#Command.Print_help)

```go
func (c *Command) Print_help()
```

打印 Help 模板

### func (\*Command) [RemoveCommand](/src/github.com/KianKw/myCobra/command.go?s=2932:2985#L140) [¶](#Command.RemoveCommand)

```go
func (cmd *Command) RemoveCommand(rmCmds ...*Command)
```

从该命令中删除子命令.

### func (\*Command) [Root](/src/github.com/KianKw/myCobra/command.go?s=4362:4397#L196) [¶](#Command.Root)

```go
func (cmd *Command) Root() *Command
```

Root 返回 root 根命令.

### func (\*Command) [Runnable](/src/github.com/KianKw/myCobra/command.go?s=1272:1307#L68) [¶](#Command.Runnable)

```go
func (cmd *Command) Runnable() bool
```

判断该命令是否可以执行

### func (\*Command) [SetHelpCommand](/src/github.com/KianKw/myCobra/command.go?s=5636:5682#L249) [¶](#Command.SetHelpCommand)

```go
func (cmd *Command) SetHelpCommand(c *Command)
```

设置 Help 命令

### func (\*Command) [SetHelpFunc](/src/github.com/KianKw/myCobra/command.go?s=4832:4891#L219) [¶](#Command.SetHelpFunc)

```go
func (cmd *Command) SetHelpFunc(f func(*Command, []string))
```

设置 Help 内容



---

---

---

编译工具：`go1.10.4`

