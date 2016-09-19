package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"

	"github.com/funny/pprof"
)

func (cmd *CMD) Shell(name string) {
	cmd.initShell(name)

	if pid := syscall.Getpid(); pid != 1 {
		ioutil.WriteFile(name+".pid", []byte(strconv.Itoa(pid)), 0777)
		defer os.Remove(name + ".pid")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGUSR1, syscall.SIGTERM)

	for sig := range sigChan {
		switch sig {
		case syscall.SIGUSR1:
			cmd.shellCommand(name)
		case syscall.SIGTERM:
			log.Println("killed")
			return
		}
	}
}

func (cmd *CMD) shellCommand(name string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("shellCommand() panic: %v\n %s", err, debug.Stack())
		}
	}()

	command, err := ioutil.ReadFile(name + ".cmd")
	if err != nil {
		log.Printf("read "+name+".cmd failed: %v", err)
		return
	}

	if _, ok := cmd.Process(string(command)); !ok {
		log.Printf("unsupported command: %s", command)
	}
}

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
		"Start or stop cpu profile. The profile will saved to "+name+".cpu.profile.",
		func(args []string) {
			switch args[1] {
			case "start":
				pprof.StartCPUProfile(name + ".cpu.profile")
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
				pprof.SaveProfile("block", name+".block.profile", debug)
			case "gc":
				log.Printf("lookup gc: %s", pprof.GCSummary())
			case "goroutine":
				debug, _ := strconv.Atoi(args[2])
				pprof.SaveProfile("goroutine", name+".goroutine.profile", debug)
			case "heap":
				debug, _ := strconv.Atoi(args[2])
				pprof.SaveProfile("heap", name+".heap.profile", debug)
			case "threadcreate":
				debug, _ := strconv.Atoi(args[2])
				pprof.SaveProfile("threadcreate", name+".threadcreate.profile", debug)
			}
		},
	)
}
