package cmd

import (
	"log"
	"os"
	"strconv"

	"github.com/funny/pprof"
)

func (cmd *CMD) initShell(name string) {
	cmd.Register(
		"help",
		"Print this screen",
		func() {
			cmd.Help(os.Stderr)
		},
	)

	cmd.Register(
		"cpuprof (start|stop)",
		"Start or stop cpu profile. The profile will saved to "+name+".cpu.prof.",
		func(args []string) {
			switch args[1] {
			case "start":
				pprof.StartCPUProfile(name + ".cpu.prof")
			case "stop":
				pprof.StopCPUProfile()
			}
		},
	)

	cmd.Register(
		"lookup (block|gc|goroutine|heap|threadcreate) ([0-9]{0,4})",
		"Dump pprof data",
		func(args []string) {
			switch args[1] {
			case "block":
				debug, _ := strconv.Atoi(args[2])
				pprof.SaveProfile("block", name+".block.prof", debug)
			case "gc":
				log.Printf("lookup gc: %s", pprof.GCSummary())
			case "goroutine":
				debug, _ := strconv.Atoi(args[2])
				pprof.SaveProfile("goroutine", name+".goroutine.prof", debug)
			case "heap":
				debug, _ := strconv.Atoi(args[2])
				pprof.SaveProfile("heap", name+".heap.prof", debug)
			case "threadcreate":
				debug, _ := strconv.Atoi(args[2])
				pprof.SaveProfile("threadcreate", name+".threadcreate.prof", debug)
			}
		},
	)
}
