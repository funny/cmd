Introduction
============

This package is used to handle text command in Go program.

Basic usage:

```go
import "github.com/funny/cmd"

cmd.Register("hello", "say hello", func() {
	fmt.Println("Hello!")
})

var command string
fmt.Scanln(&command)

cmd.Process(command)
```

The first argument of `Register()` function, is a regular expression. 

When a command match the regular expression the handler function will be invoked.

A command handler function can be any kind of these types:

```go
func()

func() interface{}

func(args []string)

func(args []string) interface{}
```

When a command handler receive `[]string` argument, the command will be splited by `Regexp.FindStringSubmatch()` and pass to handler.

For example:

```go
import "github.com/funny/cmd"

cmd.Register("cmd (abc|def) ([0-9]) ([0-9])", "test", func(args []string){
	fmt.Println(args[0]) // cmd abc 1 2
	fmt.Println(args[1]) // abc
	fmt.Println(args[2]) // 1
	fmt.Println(args[3]) // 2
})

cmd.Process("cmd abc 1 2")
```

Document: [http://godoc.org/github.com/funny/cmd](http://godoc.org/github.com/funny/cmd)