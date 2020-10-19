# tinyini 包

## 概述 

tinyini 是一个能够解析配置文件，同时根据配置文件的修改做监听与相应的包。



## 目录 

* 变量
    * type Listener
    * type ListenFunc
    * type config
    * type errorString
* watch.go 文件
    * func init
    * func listen
    * func Watch
* listen.go 文件
    * func keepLoadListener
* file.go 文件
    * func PathExists
    * func readFile
    * func loadFile
* config.go 文件
    * func isSame
    * func display
* error.go 文件



## 包内文件

* [config.go](./config.go)
* [file.go](./file.go)
* [listen.go](./listen.go)
* [watch.go](./watch.go)
* [error.go](./error.go)



## 变量

### type Listener

Listener 接口

```go
type Listener interface {
	listen(filename string)
}
```



### type ListenFunc

Listener 函数

```go
type ListenFunc func(filename string) (*config, error)
```



### type config

配置信息结构体。

```go
type config struct {
	filePath string
	info map[string]map[string]string
}
```



### type errorString

自定义错误类型

```go
type errorString struct {
	s string
}
```

定义错误信息

```go
var (
	ErrNoFile = New("no such file")
	ErrOpenFile = New("Open file failed!")
	ErrReadFile = New("Can not read the file!")
)
```



## watch.go 文件

### func init

```go
func init()
```

初始化函数，根据操作系统类型来确定注释的符号是 `;` 还是 `#`。

### func listen

```go
func (f ListenFunc) listen(filename string) (*config, error)
```

listen 函数。

### func Watch

```go
func Watch(filename string, listener ListenFunc)
```

检测函数。



## listen.go 文件

### func keepLoadListener

```go
func keepLoadListener(filename string) (*config, error)
```

持续读文件的监听方式。运行过程中不断地读取文件，发现有变化后，打印新的信息并更新内存中的配置信息。



## file.go 文件

### func PathExists

```go
func PathExists(path string) bool
```

判断文件路径是否存在，并打印错误信息。

### func readFile

```go
func readFile(filename string) (*os.File, error)
```

将文件中的内容读取到内存中。

### func loadFile

```go
func loadFile(filename string) (*config, error)
```

将内存中的文件内容加载到 config 结构体中，输入文件名，输出文件内容所构成的结构体。



## config.go 文件

### func isSame

```go
func isSame(c1 *config, c2 *config) bool
```

判断两个 config 结构体内容是否相同，如果相同则返回 true，否则返回 false。

### func display

```go
func (c *config) display()
```

打印结构体信息。



## error.go 文件

### func New

```go
func New(text string) error
```

生成错误信息，输入一个字符串，输出一个错误类型。



---

---

---

编译工具：`go1.10.4`

