package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func main() {
	pixelgl.Run(run)
}

const width = 1400
const height = 600

var maxIterations = 200

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, width, height),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	minX, maxX := -4.413230702312, 2.905435964354
	minY, maxY := -1.365239583949, 1.274760416051
	zoomMinX, zoomMaxX := 1.24417333333/8, -1.24417333333/8
	zoomMinY, zoomMaxY := 0.4488/8, -0.4488/8
	i := 0.0
	for !win.Closed() {
		i++
		pic := pixel.PictureDataFromImage(createImage(minX, maxX, minY, maxY))
		sprite := pixel.NewSprite(pic, pic.Bounds())
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		maxIterations += int(i * 2)
		zoomMinX = zoomMinX * 0.97
		zoomMaxX = zoomMaxX * 0.97
		zoomMinY = zoomMinY * 0.97
		zoomMaxY = zoomMaxY * 0.97
		minX += zoomMinX
		maxX += zoomMaxX
		minY += zoomMinY
		maxY += zoomMaxY
		fmt.Println("---------------------[START (", i, ")]----------------")
		fmt.Println("zoomMinX:", zoomMinX)
		fmt.Println("zoomMaxX:", zoomMaxX)
		fmt.Println("zoomMinY:", zoomMinY)
		fmt.Println("zoomMaxY:", zoomMaxY)

		fmt.Println("minX:", minX)
		fmt.Println("maxX:", maxX)
		fmt.Println("minY:", minY)
		fmt.Println("maxY:", maxY)
		fmt.Println("---------------------[END (", i, ")]-----------------")
		win.Update()
	}
}

func createImage(minX, maxX, minY, maxY float64) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	var wg sync.WaitGroup
	wg.Add(img.Bounds().Dy())
	for y := 0; y < img.Bounds().Dy(); y++ {
		go func(y int) {
			defer wg.Done()
			for x := 0; x < img.Bounds().Dx(); x++ {
				r := scale(0, width, minX, maxX, float64(x))
				i := scale(0, height, minY, maxY, float64(y))
				c := complex(r, i)
				n, inSet := isInSet(c)
				if inSet {
					img.Set(x, y, colornames.Black)
				} else {
					rainbowIndex := scale(0, float64(maxIterations), 0, maxRainbow, float64(n))
					rainbowIndex = math.Sqrt(rainbowIndex) * rainbowIndex
					img.Set(x, y, colorFromIndex(int(rainbowIndex)))
				}
			}
		}(y)
	}
	wg.Wait()
	return img
}

const maxRainbow = 256 * 3

func section(n int) int {
	return 256 * n
}

func colorFromIndex(i int) color.RGBA {
	i = i % section(5)
	// Red to yellow.
	if i < section(1) {
		return color.RGBA{255, uint8(i), 0, 255}
	}
	// Yellow to green.
	if i < section(2) {
		return color.RGBA{uint8(section(2) - i - 1), 255, 0, 255}
	}
	// Green to light blue.
	if i < section(3) {
		return color.RGBA{0, 255, uint8(section(2) + i), 255}
	}
	// Light blue to dark blue.
	if i < section(4) {
		return color.RGBA{0, uint8(section(4) - i - 1), 255, 255}
	}
	// Dark blue to purple.
	return color.RGBA{uint8(section(4) + i), 0, 255, 255}
}

func scale(fromMin, fromMax, toMin, toMax float64, v float64) float64 {
	return ((v / (fromMax - fromMin)) * (toMax - toMin)) + toMin
}

func isInSet(c complex128) (n int, inSet bool) {
	z := c
	for n = 0; n < maxIterations; n++ {
		z = (z * z) + c
		if real(z) > 2 || imag(z) > 2 {
			return
		}
	}
	inSet = true
	return
}

func save(fileName string, img image.Image) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	return png.Encode(f, img)
}
