package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// change this to process you care about
const PROCESS = "~PROCESS~"

func main() {
	rc, pid := system("/bin/pidof", "-s", PROCESS)
	if rc > 0 {
		exit(rc, pid)
	}
	prod_pid_fd := fmt.Sprintf("/proc/%s/fd", bytes.Trim(pid, "\n"))
	rc, out := system("/bin/ls", prod_pid_fd)
	if rc > 0 {
		exit(rc, out)
	}
	cnt := bytes.Count(out, []byte("\n"))
	fmt.Println(cnt)
}

func exit(exit int, err []byte) {
	fmt.Print(string(err))
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
