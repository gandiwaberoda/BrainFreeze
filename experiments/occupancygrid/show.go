package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/eiannone/keyboard"
	"github.com/go-p5/p5"
)

var (
	img   image.Image
	robot = NewSensorModel()

	realPosition    WorldCordinate
	realOrientation float64

	mcl = NewMonteCarlo(winW, winH)
)

const (
	winW = 720
	winH = 1020
)

// Untuk MCL
var (
	firstMclFrame       bool        = true
	lastPositionReading image.Point = image.Point{}
	lastRotationReading float64
)

func KeyboardListenerWorker() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press ESC to quit")
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		fmt.Printf("You pressed: rune %q, key %X\r\n", char, key)
		if key == keyboard.KeyEsc {
			break
		} else if key == keyboard.KeyArrowLeft {
			robot.WorldRot -= 10
			fmt.Println("Current orientation:", robot.WorldRot)
		} else if key == keyboard.KeyArrowRight {
			robot.WorldRot += 10
			fmt.Println("Current orientation:", robot.WorldRot)
		}
	}
}

func main() {
	go KeyboardListenerWorker()
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

// func drawMclParticles() {
// 	for _, v := range mcl.particles {
// 		// if v.weight == 0.0 {
// 		// 	// fmt.Println("???")
// 		// 	continue
// 		// }
// 		if !v.use {
// 			// fmt.Println(v)
// 			continue
// 		}
// 		p5.Circle(float64(v.x), float64(v.y), 3+(v.weight*500))
// 	}
// }

func drawRobotSenses(reading map[float64]LidarReading, orientation float64, pos WorldCordinate) {
	_x := 20.0
	_y := 20.0

	cc := pos.AsCanvasCordinate()
	centerX := float64(cc.X)
	centerY := float64(cc.Y)

	for k, v := range robot.Reading {
		x, y := robot.Polar2Cartesian(k, v.ClosestPoint)
		p5.Line(centerX, centerY, x+centerX, y+centerY)
	}

	p5.Push()
	p5.Translate(centerX, centerY)

	p5.Rotate(-(orientation) * math.Pi / 180)
	p5.Triangle(-_x, 0, 0, -_y, _x, 0)
	p5.Fill(color.Transparent)
	p5.Circle(0, 0, robot.deadZoneRad)

	p5.Pop()

	p5.Stroke(color.RGBA{128, 128, 255, 255})
	p5.Fill(color.Transparent)
	p5.Circle(centerX, centerY, robot.maxDistance*2)
	p5.Pop()
}

func draw() {
	p5.DrawImage(img, 0, 0)

	mouseX := p5.Event.Mouse.Position.X
	mouseY := p5.Event.Mouse.Position.Y
	// if mouseX < 0 || mouseX > winW || mouseY < 0 || mouseY > winH {
	// 	return
	// }

	p5.Push()
	p5.Stroke(color.RGBA{255, 0, 0, 255})
	p5.StrokeWidth(4)
	if p5.Event.Mouse.Pressed && p5.Event.Mouse.Buttons.Contain(p5.ButtonLeft) {
		realPosition = CanvasCordinate{int(p5.Event.Mouse.Position.X), int(p5.Event.Mouse.Position.Y)}.AsWorldCordinate()
		realOrientation = robot.WorldRot

		fmt.Println("Canvas Pos:", realPosition)
	}
	// robot.SenseFromImage(img, realPosition)
	// drawRobotSenses(robot.Reading, realOrientation, float64(realPosition.X), float64(realPosition.Y))
	p5.Pop()

	// Gambar Robot Cursor
	p5.Push()
	p5.Stroke(color.RGBA{128, 128, 255, 255})
	robot.SenseFromImage(img, image.Point{int(mouseX), int(mouseY)})
	drawRobotSenses(robot.Reading, robot.WorldRot, CanvasCordinate{int(p5.Event.Mouse.Position.X), int(p5.Event.Mouse.Position.Y)}.AsWorldCordinate())

	// MCL
	// p5.Push()
	// drawMclParticles()
	// p5.Pop()

	// if p5.Event.Mouse.Buttons.Contain(p5.ButtonLeft) {
	// 	curPos := image.Point{int(mouseX), int(mouseY)}

	// 	if firstMclFrame {
	// 		lastPositionReading = curPos
	// 		lastRotationReading = realOrientation
	// 		firstMclFrame = false
	// 	}

	// 	delta := curPos.Sub(lastPositionReading)
	// 	deltaRot := realOrientation - lastRotationReading

	// 	// fmt.Println("Haloha", delta, deltaRot, lastFrameMcl)
	// 	mcl.Update(float64(delta.X), float64(delta.Y), deltaRot, func(x, y, rot float64) float64 {
	// 		robot.WorldRot = rot
	// 		robot.SenseFromImage(img, image.Point{int(x), int(y)})
	// 		atPredictedReading := make(map[float64]LidarReading)
	// 		for k, v := range robot.Reading {
	// 			atPredictedReading[k] = v
	// 		}

	// 		robot.WorldRot = realOrientation
	// 		robot.SenseFromImage(img, realPosition)
	// 		sensorReading := robot.Reading

	// 		// fmt.Println(atPredictedReading, sensorReading)

	// 		totalErr := 0.0
	// 		for k, _ := range atPredictedReading {
	// 			err := math.Sqrt(math.Pow(atPredictedReading[k].ClosestPoint-sensorReading[k].ClosestPoint, 2))
	// 			totalErr += err
	// 		}
	// 		// fmt.Println(totalErr)
	// 		return totalErr
	// 	})
	// 	mcl.Resample()

	// 	lastPositionReading = curPos
	// 	lastRotationReading = realOrientation
	// }

	// // MCL Estimated Pose
	// {
	// 	p5.Push()
	// 	// p5.Stroke(color.RGBA{0, 0, 255, 255})
	// 	x, y, _ := mcl.EstimatePose()
	// 	// fmt.Println(x, y)
	// 	p5.Fill(color.RGBA{255, 0, 0, 255})
	// 	p5.Circle(x, y, 20)
	// 	p5.Pop()
	// }

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

type WorldCordinate image.Point
type CanvasCordinate image.Point

func (c CanvasCordinate) AsWorldCordinate() WorldCordinate {
	x := c.X - 60
	y := 1020 - 60 - c.Y // WinH
	return WorldCordinate{X: x, Y: y}
}

func (c WorldCordinate) AsCanvasCordinate() CanvasCordinate {
	x := c.X + 60
	y := 1020 - 60 - c.Y // WinH
	return CanvasCordinate{X: x, Y: y}
}
