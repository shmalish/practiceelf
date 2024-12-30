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
