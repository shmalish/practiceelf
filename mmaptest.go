package main

/*
#include <stdint.h>
#include <stdio.h>
#include <sys/mman.h>
#include <unistd.h>
#include <stdlib.h>

// We can use GCC's built-in __clear_cache on some platforms (x86 is usually no-op).
void flushInstructionCache(char *beg, char *end) {
#if defined(__GNUC__)
    __builtin___clear_cache(beg, end);
#endif
}
*/
import "C"

import (
    "fmt"
    "io"
    "os"
    "syscall"
    "unsafe"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s <shellcode_file>\n", os.Args[0])
        os.Exit(1)
    }

    // Open the shellcode file
    shellcodeFile := os.Args[1]
    f, err := os.Open(shellcodeFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to open shellcode file: %v\n", err)
        os.Exit(1)
    }
    defer f.Close()

    // Get file size
    fi, err := f.Stat()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to stat file: %v\n", err)
        os.Exit(1)
    }
    size := fi.Size()
    if size == 0 {
        fmt.Fprintf(os.Stderr, "Shellcode file is empty.\n")
        os.Exit(1)
    }
    fmt.Printf("Shellcode size: %d bytes\n", size)

    // Map memory as RW (no exec yet). Using syscall.Mmap is simpler than raw Syscall.
    data, err := syscall.Mmap(
        -1,                            // file descriptor (not used with MAP_ANONYMOUS)
        0,                             // offset
        int(size),                     // length
        syscall.PROT_READ|syscall.PROT_WRITE,
        syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS,
    )
    if err != nil {
        fmt.Fprintf(os.Stderr, "Mmap RW failed: %v\n", err)
        os.Exit(1)
    }

    // Read the shellcode into the mapped region
    _, err = io.ReadFull(f, data)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to read shellcode file: %v\n", err)
        _ = syscall.Munmap(data)
        os.Exit(1)
    }

    // Flip the memory to RX
    err = syscall.Mprotect(data, syscall.PROT_READ|syscall.PROT_EXEC)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Mprotect RX failed: %v\n", err)
        _ = syscall.Munmap(data)
        os.Exit(1)
    }

    // Flush the CPU instruction cache (mainly needed for non-x86)
    funcPtr := uintptr(unsafe.Pointer(&data[0]))
    length := uintptr(len(data))
    C.flushInstructionCache(
        (*C.char)(unsafe.Pointer(funcPtr)),
        (*C.char)(unsafe.Pointer(funcPtr+length)),
    )

    // Convert the mapped region to a "no-arg function" pointer
    sc := *(*func())(unsafe.Pointer(&data[0]))

    fmt.Println("[*] Executing shellcode...")
    // Potentially unsafe! Make sure your shellcode is correct.
    sc()

    // Unmap (not strictly necessary if we exit, but it's good practice)
    _ = syscall.Munmap(data)
}
