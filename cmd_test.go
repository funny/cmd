package cmd

import (
	"github.com/funny/utest"
	"testing"
)

func Test_Cmd(t *testing.T) {
	Register("lookup api", "save api time record to file", func(args []string) interface{} {
		utest.EqualNow(t, args[0], "lookup api")
		return 1
	})

	Register("lookup gc", "save gc summary to file", func(args []string) interface{} {
		utest.EqualNow(t, args[0], "lookup gc")
		return 2
	})

	Register("lookup heap ([0-2])", "save heap status to file", func(args []string) interface{} {
		utest.EqualNow(t, args[1], "1")
		return 3
	})

	n, ok := Process("lookup api")
	utest.Assert(t, ok && n == 1, n)
	Process("lookup gc")

	n, ok = Process("lookup gc")
	utest.Assert(t, ok && n == 2, n)

	n, ok = Process("lookup heap 1")
	utest.Assert(t, ok && n == 3, n)
}
