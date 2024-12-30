package main

import (
    "fmt"
    "io"
    "os"
    "unsafe"

    "golang.org/x/sys/unix"
)

// AT_EMPTY_PATH is the flag telling execveat to use the FD even if the path is "".
const AT_EMPTY_PATH = 0x1000

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s <hello_world_binary>\n", os.Args[0])
        os.Exit(1)
    }
    binaryPath := os.Args[1]

    // 1. Read the ELF from disk
    data, err := os.ReadFile(binaryPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to read file %s: %v\n", binaryPath, err)
        os.Exit(1)
    }

    // 2. memfd_create: an in-memory file descriptor
    fd, err := unix.MemfdCreate("memexec", 0)
    if err != nil {
        fmt.Fprintf(os.Stderr, "MemfdCreate error: %v\n", err)
        os.Exit(1)
    }

    // 3. Write ELF data into the memfd
    n, err := unix.Write(fd, data)
    if err != nil || n != len(data) {
        fmt.Fprintf(os.Stderr, "Write error: %v (wrote %d of %d)\n", err, n, len(data))
        _ = unix.Close(fd)
        os.Exit(1)
    }

    // Rewind to start
    if off, err := unix.Seek(fd, 0, io.SeekStart); err != nil || off != 0 {
        fmt.Fprintf(os.Stderr, "Seek error: %v\n", err)
        _ = unix.Close(fd)
        os.Exit(1)
    }

    // 4. execveat(fd, "", argv, envp, AT_EMPTY_PATH) replaces *this* process with the new ELF
    argv := []string{"hello_from_memfd"} // argv[0] can be anything you like
    envp := os.Environ()

    // We never return on success
    err = execveatEmptyPath(fd, argv, envp)
    if err != nil {
        fmt.Fprintf(os.Stderr, "execveat error: %v\n", err)
        _ = unix.Close(fd)
        os.Exit(1)
    }

    // not reached if execveat succeeds
    _ = unix.Close(fd)
}

// execveatEmptyPath calls execveat(fd, "", argv, envp, AT_EMPTY_PATH).
// This is supported on Linux >= 3.19. If SYS_EXECVEAT is missing, you may need to define it manually.
func execveatEmptyPath(fd int, argv []string, envp []string) error {
    // Convert Go strings -> null-terminated C strings
    argvPtrs := make([]*byte, len(argv)+1)
    for i, s := range argv {
        cstr, err := unix.BytePtrFromString(s)
        if err != nil {
            return fmt.Errorf("invalid argv string %q: %w", s, err)
        }
        argvPtrs[i] = cstr
    }

    envpPtrs := make([]*byte, len(envp)+1)
    for i, e := range envp {
        cstr, err := unix.BytePtrFromString(e)
        if err != nil {
            return fmt.Errorf("invalid env string %q: %w", e, err)
        }
        envpPtrs[i] = cstr
    }

    // On x86_64, SYS_EXECVEAT is typically 322. On ARM64, 281, etc.
    // The "golang.org/x/sys/unix" package usually defines unix.SYS_EXECVEAT for you.
    const SYS_EXECVEAT = unix.SYS_EXECVEAT

    // A pointer to an empty C-string ("")
    emptyStringPtr, _ := unix.BytePtrFromString("")

    // int execveat(int dirfd, const char *pathname,
    //              char *const argv[], char *const envp[],
    //              int flags);
    r1, _, e := unix.Syscall6(
        SYS_EXECVEAT,
        uintptr(fd),
        uintptr(unsafe.Pointer(emptyStringPtr)),
        uintptr(unsafe.Pointer(&argvPtrs[0])),
        uintptr(unsafe.Pointer(&envpPtrs[0])),
        uintptr(AT_EMPTY_PATH),
        0,
    )
    if e != 0 {
        return e
    }
    if r1 != 0 {
        return fmt.Errorf("execveat returned %d unexpectedly", r1)
    }
    // execveat does not return on success
    return nil
}
