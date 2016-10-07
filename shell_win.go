// +build windows

package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"
)

func (cmd *CMD) Shell(name string) {
	cmd.initShell(name)

	if pid := syscall.Getpid(); pid != 1 {
		ioutil.WriteFile(name+".pid", []byte(strconv.Itoa(pid)), 0777)
		defer os.Remove(name + ".pid")
	}

	ticker := time.NewTicker(time.Second * 3)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case <-ticker.C:
			cmd.shellCommand(name)
		case <-sigChan:
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

	if _, err := os.Stat(name + ".cmd"); err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Printf("check "+name+".cmd exists failed: %v", err)
	}

	command, err := ioutil.ReadFile(name + ".cmd")
	if err != nil {
		log.Printf("read "+name+".cmd failed: %v", err)
		return
	}

	if err := os.Remove(name + ".cmd"); err != nil {
		log.Printf("remove "+name+".cmd failed: %v", err)
	}

	if _, ok := cmd.Process(string(command)); !ok {
		log.Printf("unsupported command: %s", command)
	}
}
