// +build !windows

package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
)

func (cmd *CMD) Shell(name string) {
	cmd.initShell(name)

	if pid := syscall.Getpid(); pid != 1 {
		ioutil.WriteFile(name+".pid", []byte(strconv.Itoa(pid)), 0777)
		defer os.Remove(name + ".pid")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGUSR1, syscall.SIGTERM, syscall.SIGINT)

	for sig := range sigChan {
		switch sig {
		case syscall.SIGUSR1:
			cmd.shellCommand(name)
		case syscall.SIGTERM, syscall.SIGINT:
			log.Println(name, "killed")
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
