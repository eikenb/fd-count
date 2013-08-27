package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

var warn int
var crit int

func init() {
	flag.IntVar(&warn, "w", 500, "# of file descriptors to trigger warning")
	flag.IntVar(&crit, "c", 900, "# of file descriptors to trigger critical")
	flag.Parse()
}

// change this to process you care about
const PROCESS = "~PROCESS~"

func main() {
	flag.Parse()
	rc, pid := system("/bin/pidof", "-s", PROCESS)
	if rc > 0 {
		// 0 as this app is not intended to check whether the app is running
		exit(0, PROCESS+": process not found\n")
	}
	prod_pid_fd := fmt.Sprintf("/proc/%s/fd", bytes.Trim(pid, "\n"))
	rc, out := system("/bin/ls", prod_pid_fd)
	if rc > 0 {
		exit(rc, string(out))
	}
	cnt := bytes.Count(out, []byte("\n"))
	msg := fmt.Sprintf("%s; file descriptor count: %d\n", PROCESS, cnt)
	switch {
	case cnt > crit:
		exit(2, "CRIT: "+msg)
	case cnt > warn:
		exit(1, "WARN: "+msg)
	}
	fmt.Print(msg)
}

func exit(exit int, err string) {
	fmt.Print(err)
	os.Exit(exit)
}

func system(name string, params ...string) (int, []byte) {
	exitstatus := 0
	cmd := exec.Command(name, params...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		exitstatus = 1
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitstatus = status.ExitStatus()
			}
		}
	}
	return exitstatus, out
}
