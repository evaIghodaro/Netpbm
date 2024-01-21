package Netpbm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// PBM represents a PBM image
type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// ReadPBM reads a PBM image from a file and returns a structure representing the image.
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Read the magic number
	magicNumber, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading magic number: %v", err)
	}
	magicNumber = strings.TrimSpace(magicNumber)
	if magicNumber != "P1" && magicNumber != "P4" {
		return nil, fmt.Errorf("invalid magic number: %s", magicNumber)
	}

	// Read dimensions
	dimensions, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading dimensions: %v", err)
	}
	var width, height int
	_, err = fmt.Sscanf(strings.TrimSpace(dimensions), "%d %d", &width, &height)
	if err != nil {
		return nil, fmt.Errorf("invalid dimensions: %v", err)
	}

	data := make([][]bool, height)

	for i := range data {
		data[i] = make([]bool, width)
	}

	if magicNumber == "P1" {
		// Read format P1 (ASCII)
		for y := 0; y < height; y++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("error reading data at line %d: %v", y, err)
			}
			fields := strings.Fields(line)
			for x, field := range fields {
				if x >= width {
					return nil, fmt.Errorf("index out of range at line %d", y)
				}
				data[y][x] = field == "1"
			}
		}
	} else if magicNumber == "P4" {
		// Read format P4 (binary)
		expectedBytesPerRow := (width + 7) / 8
		for y := 0; y < height; y++ {
			row := make([]byte, expectedBytesPerRow)
			n, err := reader.Read(row)
			if err != nil {
				if err == io.EOF {
					return nil, fmt.Errorf("unexpected end of file at line %d", y)
				}
				return nil, fmt.Errorf("error reading pixel data at line %d: %v", y, err)
			}
			if n < expectedBytesPerRow {
				return nil, fmt.Errorf("unexpected end of file at line %d, expected %d bytes, got %d", y, expectedBytesPerRow, n)
			}

			for x := 0; x < width; x++ {
				byteIndex := x / 8
				bitIndex := 7 - (x % 8)

				// Convert ASCII to decimal and extract the bit
				decimalValue := int(row[byteIndex])
				bitValue := (decimalValue >> bitIndex) & 1

				data[y][x] = bitValue != 0
			}
		}
	}

	return &PBM{data, width, height, magicNumber}, nil
}

// Save saves a PBM image to a file.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write the magic number
	_, err = fmt.Fprintln(writer, pbm.magicNumber)
	if err != nil {
		return fmt.Errorf("error writing magic number: %v", err)
	}

	// Write dimensions
	_, err = fmt.Fprintf(writer, "%d %d\n", pbm.width, pbm.height)
	if err != nil {
		return fmt.Errorf("error writing dimensions: %v", err)
	}

	if pbm.magicNumber == "P1" {
		// Write format P1 (ASCII)
		for y := 0; y < pbm.height; y++ {
			for x := 0; x < pbm.width; x++ {
				if pbm.data[y][x] {
					_, err := fmt.Fprint(writer, "1 ")
					if err != nil {
						return fmt.Errorf("error writing data at line %d, column %d: %v", y, x, err)
					}
				} else {
					_, err := fmt.Fprint(writer, "0 ")
					if err != nil {
						return fmt.Errorf("error writing data at line %d, column %d: %v", y, x, err)
					}
				}
			}
			_, err := fmt.Fprintln(writer, "")
			if err != nil {
				return fmt.Errorf("error writing newline: %v", err)
			}
		}
	} else if pbm.magicNumber == "P4" {
		// Write format P4 (binary)
		for y := 0; y < pbm.height; y++ {
			var currentByte byte
			for x := 0; x < pbm.width; x++ {
				bitIndex := 7 - (x % 8)
				bitValue := 0
				if pbm.data[y][x] {
					bitValue = 1
				}
				// Update the appropriate bit in the byte
				currentByte |= byte(bitValue << bitIndex)

				if (x+1)%8 == 0 || x == pbm.width-1 {
					_, err := writer.Write([]byte{currentByte})
					if err != nil {
						return fmt.Errorf("error writing binary data at line %d: %v", y, err)
					}
					currentByte = 0
				}
			}
		}
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing write buffer: %v", err)
	}

	return nil
}

// Size returns the width and height of the PBM image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At returns the pixel value at position (x, y) in the PBM image.
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

// Set sets the pixel value at position (x, y) in the PBM image.
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

// Invert inverts the values of all pixels in the PBM image.
func (pbm *PBM) Invert() {
	for y := 0; y < pbm.height; y++ {
		for x := 0; x < pbm.width; x++ {
			pbm.data[y][x] = !pbm.data[y][x]
		}
	}
}

// Flip flips the PBM image horizontally.
func (pbm *PBM) Flip() {
	for y := 0; y < pbm.height; y++ {
		for x := 0; x < pbm.width/2; x++ {
			pbm.data[y][x], pbm.data[y][pbm.width-x-1] = pbm.data[y][pbm.width-x-1], pbm.data[y][x]
		}
	}
}

// Flop flips the PBM image vertically.
func (pbm *PBM) Flop() {
	for x := 0; x < pbm.width; x++ {
		for y := 0; y < pbm.height/2; y++ {
			pbm.data[y][x], pbm.data[pbm.height-y-1][x] = pbm.data[pbm.height-y-1][x], pbm.data[y][x]
		}
	}
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
