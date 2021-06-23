package models

import "fmt"

type Force struct {
	x        float64
	y        float64
	rot      float64
	handling float64
}

func (f *Force) Idle() {
	f.x = 0
	f.y = 0
	f.rot = 0
}

//#region
func (f *Force) AddX(_x float64) {
	f.x += _x
}
func (f *Force) AddY(_y float64) {
	f.y += _y
}

func (f Force) GetX() float64 {
	return f.x
}
func (f *Force) AddRot(_rot Degree) {
	f.rot += float64(_rot)
}
func (f *Force) HandlingReverse() {
	f.handling = -1.0
}

func (f *Force) EnableHandling() {
	f.handling = 1.0
}

func (f *Force) DisableHandling() {
	f.handling = 0.0
}

//#region
func (f Force) GetY() float64 {
	return f.y
}

func (f Force) GetRot() float64 {
	return f.rot
}

func (f Force) GetHandling() float64 {
	return f.handling
}

//#region
func (f Force) AsGutCommandString() string {
	return fmt.Sprint("*a,", f.GetX(), ",", f.GetY(), ",", f.GetRot(), ",", f.GetHandling(), "#")
}
