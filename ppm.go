package Netpbm

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

// PPM structure represents a Portable Pixmap image
type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           uint8
}

// Pixel structure represents a single pixel with RGB values
type Pixel struct {
	R, G, B uint8
}

// Point structure represents a 2D point
type Point struct {
	X, Y int
}

// ReadPPM reads a PPM image from the specified file name
func ReadPPM(filename string) (*PPM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	ppm := &PPM{}

	// Read the magic number
	scanner.Scan()
	ppm.magicNumber = scanner.Text()

	// Read width and height
	scanner.Scan()
	ppm.width, _ = strconv.Atoi(scanner.Text())
	scanner.Scan()
	ppm.height, _ = strconv.Atoi(scanner.Text())

	// Read the maximum pixel value
	scanner.Scan()
	maxValue, _ := strconv.Atoi(scanner.Text())
	ppm.max = uint8(maxValue)

	// Initialize the data slice
	ppm.data = make([][]Pixel, ppm.height)
	for i := range ppm.data {
		ppm.data[i] = make([]Pixel, ppm.width)
	}

	// Read pixel values
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			scanner.Scan()
			ppm.data[i][j].R, _ = strconv.ParseUint(scanner.Text(), 10, 8)
			scanner.Scan()
			ppm.data[i][j].G, _ = strconv.ParseUint(scanner.Text(), 10, 8)
			scanner.Scan()
			ppm.data[i][j].B, _ = strconv.ParseUint(scanner.Text(), 10, 8)
		}
	}

	return ppm, nil
}

// Size returns the width and height of the PPM image
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// At returns the pixel value at the specified coordinates (x, y)
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set updates the pixel value at the specified coordinates (x, y)
func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

// Save writes the PPM image to the specified file
func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Write magic number, width, height, and maximum pixel value
	fmt.Fprintf(writer, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)

	// Write pixel values
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			fmt.Fprintf(writer, "%d %d %d ", ppm.data[i][j].R, ppm.data[i][j].G, ppm.data[i][j].B)
		}
		fmt.Fprintln(writer)
	}

	return nil
}

// Invert inverts the colors of the PPM image
func (ppm *PPM) Invert() {
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			ppm.data[i][j].R = ppm.max - ppm.data[i][j].R
			ppm.data[i][j].G = ppm.max - ppm.data[i][j].G
			ppm.data[i][j].B = ppm.max - ppm.data[i][j].B
		}
	}
}

// Flip flips the PPM image horizontally
func (ppm *PPM) Flip() {
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width/2; j++ {
			ppm.data[i][j], ppm.data[i][ppm.width-1-j] = ppm.data[i][ppm.width-1-j], ppm.data[i][j]
		}
	}
}

// Flop flips the PPM image vertically
func (ppm *PPM) Flop() {
	for i := 0; i < ppm.height/2; i++ {
		ppm.data[i], ppm.data[ppm.height-1-i] = ppm.data[ppm.height-1-i], ppm.data[i]
	}
}

// SetMagicNumber sets the magic number of the PPM image
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue sets the maximum pixel value of the PPM image
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	ppm.max = maxValue
}

// Rotate90CW rotates the PPM image 90 degrees clockwise
func (ppm *PPM) Rotate90CW() {
	newData := make([][]Pixel, ppm.width)
	for i := range newData {
		newData[i] = make([]Pixel, ppm.height)
	}
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			newData[j][ppm.height-1-i] = ppm.data[i][j]
		}
	}
	ppm.width, ppm.height = ppm.height, ppm.width
	ppm.data = newData
}

// ToPGM converts the PPM image to a PGM image (grayscale)
func (ppm *PPM) ToPGM() *PGM {
	pgm := &PGM{
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P2",
		max:         ppm.max,
		data:        make([][]uint8, ppm.height),
	}
	for i := range pgm.data {
		pgm.data[i] = make([]uint8, ppm.width)
		for j := 0; j < ppm.width; j++ {
			// Convert RGB to grayscale using the luminosity formula
			pgm.data[i][j] = uint8(0.299*float64(ppm.data[i][j].R) + 0.587*float64(ppm.data[i][j].G) + 0.114*float64(ppm.data[i][j].B))
		}
	}
	return pgm
}

// ToPBM converts the PPM image to a PBM image (black and white)
func (ppm *PPM) ToPBM() *PBM {
	pbm := &PBM{
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P1",
		data:        make([][]bool, ppm.height),
	}
	for i := range pbm.data {
		pbm.data[i] = make([]bool, ppm.width)
		for j := 0; j < ppm.width; j++ {
			// Convert RGB to binary using a simple threshold (128)
			grayValue := 0.299*float64(ppm.data[i][j].R) + 0.587*float64(ppm.data[i][j].G) + 0.114*float64(ppm.data[i][j].B)
			pbm.data[i][j] = grayValue > 128
		}
	}
	return pbm
}

// DrawLine draws a line on the PPM image between two points with the specified color
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	// Implement the DrawLine function here
	deltaX := p2.X - p1.X
	deltaY := p2.Y - p1.Y
	steps := int(math.Max(math.Abs(float64(deltaX)), math.Abs(float64(deltaY))))
	xIncrement := float64(deltaX) / float64(steps)
	yIncrement := float64(deltaY) / float64(steps)
	x := float64(p1.X)
	y := float64(p1.Y)
	for i := 0; i <= steps; i++ {
		ppm.Set(int(x), int(y), color)
		x += xIncrement
		y += yIncrement
	}
}

// DrawRectangle draws a rectangle on the PPM image with the specified color
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	p2 := Point{p1.X + width, p1.Y}
	p3 := Point{p1.X + width, p1.Y + height}
	p4 := Point{p1.X, p1.Y + height}
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p4, color)
	ppm.DrawLine(p4, p1, color)
}

// DrawFilledRectangle draws a filled rectangle on the PPM image with the specified color
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			ppm.Set(p1.X+j, p1.Y+i, color)
		}
	}
}

// DrawCircle draws a circle on the PPM image with the specified color
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	for x := -radius; x <= radius; x++ {
		for y := -radius; y <= radius; y++ {
			if x*x+y*y <= radius*radius {
				ppm.Set(center.X+x, center.Y+y, color)
			}
		}
	}
}

// DrawFilledCircle draws a filled circle on the PPM image with the specified color
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	for x := -radius; x <= radius; x++ {
		for y := -radius; y <= radius; y++ {
			if x*x+y*y <= radius*radius {
				ppm.Set(center.X+x, center.Y+y, color)
			}
		}
	}
}

// DrawTriangle draws a triangle on the PPM image with the specified color
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

// DrawFilledTriangle draws a filled triangle on the PPM image with the specified color
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	vertices := []Point{p1, p2, p3}
	sort.Slice(vertices, func(i, j int) bool {
		return vertices[i].Y < vertices[j].Y
	})
	slope1 := float64(vertices[1].X-vertices[0].X) / float64(vertices[1].Y-vertices[0].Y)
	slope2 := float64(vertices[2].X-vertices[0].X) / float64(vertices[2].Y-vertices[0].Y)
	x1 := float64(vertices[0].X)
	x2 := float64(vertices[0].X)
	for y := vertices[0].Y; y <= vertices[1].Y; y++ {
		for x := int(math.Min(x1, x2)); x <= int(math.Max(x1, x2)); x++ {
			ppm.Set(x, y, color)
		}
		x1 += slope1
		x2 += slope2
	}
	slope3 := float64(vertices[2].X-vertices[1].X) / float64(vertices[2].Y-vertices[1].Y)
	x1 = float64(vertices[1].X)
	for y := vertices[1].Y + 1; y <= vertices[2].Y; y++ {
		for x := int(math.Min(x1, x2)); x <= int(math.Max(x1, x2)); x++ {
			ppm.Set(x, y, color)
		}
		x1 += slope3
		x2 += slope2
	}
}

// DrawPolygon draws a polygon on the PPM image with the specified color
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	for i := 0; i < len(points); i++ {
		p1 := points[i]
		p2 := points[(i+1)%len(points)]
		ppm.DrawLine(p1, p2, color)
	}
}

// DrawFilledPolygon draws a filled polygon on the PPM image with the specified color
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	// Find the bounding edges of the polygon
	minX, minY := points[0].X, points[0].Y
	maxX, maxY := points[0].X, points[0].Y
	for _, p := range points {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	// Create an array to store intersections per row
	intersections := make([]int, maxY-minY+1)
	// Iterate over each edge of the polygon
	for i := 0; i < len(points); i++ {
		p1 := points[i]
		p2 := points[(i+1)%len(points)]

		// Find the minimum and maximum y-coordinates of the edge
		yMin := int(math.Min(float64(p1.Y), float64(p2.Y)))
		yMax := int(math.Max(float64(p1.Y), float64(p2.Y)))

		// Skip horizontal edges
		if yMin == yMax {
			continue
		}

		// Iterate over each row the edge crosses and update intersections
		for y := yMin; y <= yMax; y++ {
			// Calculate x-coordinate of the intersection
			xIntersection := int(float64(p1.X) + float64(y-yMin)*(float64(p2.X)-float64(p1.X))/(float64(p2.Y)-float64(p1.Y)))

			// Increment the intersection count for the current row
			intersections[y-minY] = xIntersection
		}
	}

	// Fill the polygon by connecting intersections on each row
	for y := 0; y <= maxY-minY; y++ {
		// Skip rows with no intersections
		if intersections[y] == 0 {
			continue
		}

		// Connect intersections on the current row
		for x := intersections[y]; x <= maxX-minX; x++ {
			ppm.Set(x+minX, y+minY, color)
		}
	}

	return ppm
}
