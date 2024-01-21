package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PGM represents a PGM image.
type PGM struct {
	data        [][]uint8 // Pixel values of the image
	width       int       // Width of the image
	height      int       // Height of the image
	magicNumber string    // PGM file format identifier
	max         uint      // Maximum pixel value (usually 255 for 8-bit PGM)
}

// ReadPGM reads a PGM image from a file and returns a structure representing the image.
func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read the magic number
	scanner.Scan()
	magicNumber := strings.TrimSpace(scanner.Text())

	// Read width, height, and maximum pixel value
	scanner.Scan()
	sizeLine := strings.Split(strings.TrimSpace(scanner.Text()), " ")
	width, err := strconv.Atoi(sizeLine[0])
	if err != nil {
		return nil, err
	}
	height, err := strconv.Atoi(sizeLine[1])
	if err != nil {
		return nil, err
	}

	scanner.Scan()
	maxValue, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return nil, err
	}

	data := make([][]uint8, height)
	for i := 0; i < height; i++ {
		data[i] = make([]uint8, width)
		scanner.Scan()
		row := strings.Fields(scanner.Text())
		for j := 0; j < width; j++ {
			value, err := strconv.ParseUint(row[j], 10, 8)
			if err != nil {
				return nil, err
			}
			data[i][j] = uint8(value)
		}
	}

	return &PGM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         uint(maxValue),
	}, nil
}

// Size returns the width and height of the image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At returns the pixel value at position (x, y).
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// Set sets the pixel value at position (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Save saves the PGM image to a file and returns an error if any.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			fmt.Fprintf(file, "%d ", pgm.data[i][j])
		}
		fmt.Fprintln(file)
	}

	return nil
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width/2; j++ {
			pgm.data[i][j], pgm.data[i][pgm.width-j-1] = pgm.data[i][pgm.width-j-1], pgm.data[i][j]
		}
	}
}

// Flop flips the PGM image vertically.
func (pgm *PGM) Flop() {
	for i := 0; i < pgm.height/2; i++ {
		for j := 0; j < pgm.width; j++ {
			pgm.data[i][j], pgm.data[pgm.height-i-1][j] = pgm.data[pgm.height-i-1][j], pgm.data[i][j]
		}
	}
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the maximum value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint) {
	pgm.max = maxValue
}

// Rotate90CW rotates the PGM image 90 degrees clockwise.
func (pgm *PGM) Rotate90CW() {
	rotatedData := make([][]uint8, pgm.width)
	for i := 0; i < pgm.width; i++ {
		rotatedData[i] = make([]uint8, pgm.height)
	}

	for i := 0; i < pgm.width; i++ {
		for j := 0; j < pgm.height; j++ {
			rotatedData[i][j] = pgm.data[pgm.height-j-1][i]
		}
	}

	pgm.data = rotatedData
	pgm.width, pgm.height = pgm.height, pgm.width
}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM {
	pbmData := make([][]bool, pgm.height)
	for i := 0; i < pgm.height; i++ {
		pbmData[i] = make([]bool, pgm.width)
		for j := 0; j < pgm.width; j++ {
			pbmData[i][j] = uint16(pgm.data[i][j]) > uint16(pgm.max)/2
		}
	}

	return &PBM{
		data:        pbmData,
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P4",
	}
}
