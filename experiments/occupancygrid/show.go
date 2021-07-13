package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/go-p5/p5"
)

var (
	img   image.Image
	robot = NewSensorModel()
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

func draw() {
	p5.DrawImage(img, 0, 0)
	mouseX := p5.Event.Mouse.Position.X
	mouseY := p5.Event.Mouse.Position.Y
	if mouseX < 0 || mouseX > winW || mouseY < 0 || mouseY > winH {
		return
	}

	p5.Push()
	robot.SenseFromImage(img, image.Point{int(mouseX), int(mouseY)})
	for k, v := range robot.Reading {
		p5.Stroke(color.RGBA{128, 128, 255, 255})
		x, y := Polar2Cartesian(k, v.ClosestPoint)
		p5.Line(mouseX, mouseY, x+mouseX, y+mouseY)
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
