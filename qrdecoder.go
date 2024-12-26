package main

import (
	"encoding/base64"
	"fmt"
	"image"
	_ "image/png" // Import PNG decoder
	"io/ioutil"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

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

func main() {
	// Decode the first QR code
	decodedPart1, err := decodeQRCode("infect_1.png")
	if err != nil {
		fmt.Println("Error decoding QR code part 1:", err)
		return
	}
	fmt.Println("Decoded Part 1:", decodedPart1)

	// Decode the second QR code
	decodedPart2, err := decodeQRCode("infect_2.png")
	if err != nil {
		fmt.Println("Error decoding QR code part 2:", err)
		return
	}
	fmt.Println("Decoded Part 2:", decodedPart2)

	decodedPart3, err := decodeQRCode("infect_3.png")
	if err != nil {
		fmt.Println("Error decoding QR code part 3:", err)
		return
	}
	fmt.Println("Decoded Part 3:", decodedPart3)
	decodedPart4, err := decodeQRCode("infect_4.png")
	if err != nil {
		fmt.Println("Error decoding QR code part 4:", err)
		return
	}
	fmt.Println("Decoded Part 3:", decodedPart4)
	decodedPart5, err := decodeQRCode("infect_5.png")
	if err != nil {
		fmt.Println("Error decoding QR code part 3:", err)
		return
	}
	fmt.Println("Decoded Part 3:", decodedPart5)

	// Combine both decoded parts (base64 encoded data)
	combinedBase64 := decodedPart1 + decodedPart2 + decodedPart3 + decodedPart4 + decodedPart5

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
