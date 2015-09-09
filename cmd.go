package cmd

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var defaultCmd = New()

func Register(format, desc string, callback interface{}) {
	defaultCmd.Register(format, desc, callback)
}

func Process(command string) (interface{}, bool) {
	return defaultCmd.Process(command)
}

func Help(w io.Writer) error {
	return defaultCmd.Help(w)
}

type handler struct {
	Format   string
	Desc     string
	Callback interface{}
	Regexp   *regexp.Regexp
}

type CMD struct {
	handlers []handler
}

func New() *CMD {
	return &CMD{}
}

func (cmd *CMD) Register(format, desc string, callback interface{}) {
	format = strings.Trim(format, "\n\r\t ")
	regexpStr := format
	if regexpStr[0] != '^' {
		regexpStr = "^" + regexpStr
	}
	if regexpStr[len(regexpStr)-1] != '$' {
		regexpStr = regexpStr + "$"
	}
	cmd.handlers = append(cmd.handlers, handler{
		format,
		desc,
		callback,
		regexp.MustCompile(regexpStr),
	})
}

func (cmd *CMD) Process(command string) (result interface{}, ok bool) {
	command = strings.Trim(command, "\n\r\t ")
	for i := 0; i < len(cmd.handlers); i++ {
		if matches := cmd.handlers[i].Regexp.FindStringSubmatch(command); len(matches) > 0 {
			switch callback := cmd.handlers[i].Callback.(type) {
			case func():
				callback()
			case func([]string):
				callback(matches)
			case func() interface{}:
				result = callback()
			case func([]string) interface{}:
				result = callback(matches)
			}
			return result, true
		}
	}
	return
}

func (cmd *CMD) Help(w io.Writer) error {
	var maxLen int
	for i := 0; i < len(cmd.handlers); i++ {
		if l := len(cmd.handlers[i].Format); l > maxLen {
			maxLen = l
		}
	}

	var fmtStr = "%-" + strconv.Itoa(maxLen) + "s\t%s\n"
	for i := 0; i < len(cmd.handlers); i++ {
		if _, err := fmt.Fprintf(w, fmtStr, cmd.handlers[i].Format, cmd.handlers[i].Desc); err != nil {
			return err
		}
	}
	return nil
}
