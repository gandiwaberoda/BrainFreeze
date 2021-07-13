package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/go-p5/p5"
)

var (
	img image.Image
)

const (
	winW = 720
	winH = 1020
)

func main() {
	p5.Run(setup, draw)
}

func setup() {
	p5.Canvas(winW, winH)
	p5.Background(color.Gray{Y: 220})

	imgLd, err := p5.ReadImage("./assets/map/Arena.png")
	if err != nil {
		fmt.Println(err)
	}
	img = imgLd

}

func Polar2Cartesian(deg, rad float64) (x, y float64) {
	radian := ((deg * -1) - 90) * math.Pi / 180
	x = rad * math.Cos(radian)
	y = rad * math.Sin(radian)
	return
}

func draw() {
	fov := 140.0
	fovStep := 15.0
	rad := 300.0

	p5.DrawImage(img, 0, 0)
	mouseX := p5.Event.Mouse.Position.X
	mouseY := p5.Event.Mouse.Position.Y
	if mouseX < 0 || mouseX > winW || mouseY < 0 || mouseY > winH {
		return
	}

	// fmt.Println("X:", int(mouseX), "Y:", int(mouseY), "\t", px, "\t")

	p5.Push()
	fmt.Println("Gambar")
	for i := -fov; i <= fov; i += fovStep {
		var endX, endY float64
		lX, lY := Polar2Cartesian(i, rad)
		lX += mouseX
		lY += mouseY
		p5.StrokeWidth(1)
		p5.Stroke(color.RGBA{50, 50, 50, 100})
		p5.Line(mouseX, mouseY, lX, lY)

		p5.StrokeWidth(5)
		p5.Stroke(color.RGBA{128, 128, 255, 255})
		for r := 10.0; r <= rad; r++ {
			endX, endY = Polar2Cartesian(i, r)
			endX += mouseX
			endY += mouseY

			px := img.At(int(endX), int(endY))
			// fmt.Println(px.RGBA())
			if IsWhite(px) {
				// p5.Line(mouseX, mouseY, endX+rand.Float64()*30, endY+rand.Float64()*30)
				p5.Line(mouseX, mouseY, endX, endY)
				break
			}
		}
		// fmt.Println(endX, "\t", endY)
	}
	p5.Pop()

	// p5.StrokeWidth(2)
	// p5.Fill(color.RGBA{R: 255, A: 208})
	// p5.Ellipse(50, 50, 80, 80)

	// p5.Fill(color.RGBA{B: 255, A: 208})
	// p5.Quad(50, 50, 80, 50, 80, 120, 60, 120)

	// p5.Fill(color.RGBA{G: 255, A: 208})
	// p5.Rect(200, 200, 50, 100)

	// p5.Fill(color.RGBA{G: 255, A: 208})
	// p5.Triangle(100, 100, 120, 120, 80, 120)

	// p5.TextSize(24)
	// p5.Text("Hello, World!", 10, 300)

	// p5.Stroke(color.Black)
	// p5.StrokeWidth(5)
	// p5.Arc(300, 100, 80, 20, 0, 1.5*math.Pi)
}

func IsWhite(ac color.Color) bool {
	max := uint32(65535)
	r, g, b, _ := ac.RGBA()
	return r == max && g == max && b == max
}

// type LidarReading struct {
// 	SensorRot    float64
// 	ClosestPoint float64
// }
// type SensorModel struct {
// 	NumSensor   int
// 	MinRotation float64
// 	MaxRotation float64
// 	Reading     []LidarReading
// }
