package main

import (

	"encoding/base64"
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

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

	// Decode the QR code
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", err
	}

	return result.GetText(), nil
}

func main() {
	
	decodedPart1, err := decodeQRCode("basic_qrcode_part_1.png")
	if err != nil {
		fmt.Println("Error decoding QR code part 1:", err)
		return
	}
	fmt.Println("Decoded Part 1:", decodedPart1)

	decodedPart2, err := decodeQRCode("basic_qrcode_part_2.png")
	if err != nil {
		fmt.Println("Error decoding QR code part 2:", err)
		return
	}
	fmt.Println("Decoded Part 2:", decodedPart2)


	combinedBase64 := decodedPart1 + decodedPart2


	decodedBytes, err := base64.StdEncoding.DecodeString(combinedBase64)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return
	}


	fmt.Println("Decoded Binary Data (first 100 bytes):", decodedBytes[:100])


	err = ioutil.WriteFile("output.bin", decodedBytes, 0644)
	if err != nil {
		fmt.Println("Error writing binary data to file:", err)
		return
	}

	fmt.Println("Binary data successfully written to 'output.bin'")
}
