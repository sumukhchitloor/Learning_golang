package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	seccomp "github.com/seccomp/libseccomp-golang"
)

func main() {
	// Define the table flag
	tableFlag := flag.Bool("table", false, "Show only syscall statistics table")

	// Parse flags
	flag.Parse()

	// Check if we have any arguments left after flags
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: ./main [--table] <program> [args...]")
		os.Exit(1)
	}

	// Initialize syscall counter
	sc := NewSyscallCounter()

	fmt.Printf("Run %v\n", args)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace: true,
	}

	cmd.Start()
	err := cmd.Wait()
	if err != nil {
		fmt.Printf("Wait returned: %v\n", err)
	}

	pid := cmd.Process.Pid
	exit := true

	var regs syscall.PtraceRegs
	for {
		if exit {
			err = syscall.PtraceGetRegs(pid, &regs)
			if err != nil {
				break
			}

			// Get syscall number for current architecture
			var syscallNum uint64
			nativeArch, err := seccomp.GetNativeArch()
			if err != nil {
				panic(err)
			}

			// ARM64 stores syscall number in X8 register (index 8)
			if nativeArch == seccomp.ArchARM64 {
				syscallNum = regs.Regs[8] // ARM64: X8 register
			} else {
				// This code won't work on ARM64, but keeping for other architectures
				fmt.Printf("Note: Running on architecture: %v, register dump: %+v\n", nativeArch, regs)
				panic("Only ARM64 is supported in this code")
			}

			// Get syscall name
			name := sc.GetName(syscallNum)

			// Print syscall name if not in table-only mode
			if !*tableFlag {
				fmt.Printf("%s\n", name)
			}

			sc.Inc(syscallNum)
		}

		err = syscall.PtraceSyscall(pid, 0)
		if err != nil {
			break
		}

		_, err = syscall.Wait4(pid, nil, 0, nil)
		if err != nil {
			break
		}

		exit = !exit
	}

	// Print syscall statistics table if requested
	if *tableFlag {
		sc.Print()
	}
}
