package cmd

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var handlers []handler

type handler struct {
	Format   string
	Desc     string
	Callback interface{}
	Regexp   *regexp.Regexp
}

func Register(format, desc string, callback interface{}) {
	format = strings.Trim(format, "\n\r\t ")
	regexpStr := format
	if regexpStr[0] != '^' {
		regexpStr = "^" + regexpStr
	}
	if regexpStr[len(regexpStr)-1] != '$' {
		regexpStr = regexpStr + "$"
	}
	handlers = append(handlers, handler{
		format,
		desc,
		callback,
		regexp.MustCompile(regexpStr),
	})
}

func Process(command string) (result interface{}, ok bool) {
	command = strings.Trim(command, "\n\r\t ")
	for i := 0; i < len(handlers); i++ {
		if matches := handlers[i].Regexp.FindStringSubmatch(command); len(matches) > 0 {
			switch callback := handlers[i].Callback.(type) {
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

func Help(w io.Writer) error {
	var maxLen int
	for i := 0; i < len(handlers); i++ {
		if l := len(handlers[i].Format); l > maxLen {
			maxLen = l
		}
	}

	var fmtStr = "%-" + strconv.Itoa(maxLen) + "s\t%s\n"
	for i := 0; i < len(handlers); i++ {
		if _, err := fmt.Fprintf(w, fmtStr, handlers[i].Format, handlers[i].Desc); err != nil {
			return err
		}
	}
	return nil
}
