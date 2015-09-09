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

Document: [http://godoc.org/github.com/funny/cmd](http://godoc.org/github.com/funny/cmd)