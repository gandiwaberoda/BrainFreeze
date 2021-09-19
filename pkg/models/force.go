package models

import (
	"fmt"
)

type Force struct {
	x        float64
	y        float64
	rot      float64
	kick     float64
	handling float64
	reset    int

	handlingHaveChanged bool
}

func (f *Force) Idle() {
	f.x = 0
	f.y = 0
	f.rot = 0
	f.kick = 0
	f.handling = 0
	f.handlingHaveChanged = true
}

//#region
func (f *Force) AddX(_x float64) {
	f.x += _x
}
func (f *Force) AddY(_y float64) {
	f.y += _y
}

func (f *Force) AddRot(_rot Degree) {
	f.rot += float64(_rot)
}
func (f *Force) HandlingReverse() {
	f.handlingHaveChanged = true
	f.handling = -1.0
}

func (f *Force) EnableHandling() {
	f.handlingHaveChanged = true
	f.handling = 1.0
}

func (f *Force) DisableHandling() {
	f.handlingHaveChanged = true
	f.handling = 0.0
}

func (f *Force) DoReset() {
	f.reset = 1
}

func (f *Force) UndoReset() {
	f.reset = 0
}

func (f *Force) HandlingHaveChanged() bool {
	return f.handlingHaveChanged
}

func (f *Force) Kick() {
	f.kick = 1.0
}

func (f *Force) CancelKick() {
	f.kick = 0.0
}

//#region
func (f Force) GetX() float64 {
	return f.x
}

func (f Force) GetY() float64 {
	return f.y
}

func (f Force) GetRot() float64 {
	return f.rot
}

func (f Force) GetKick() float64 {
	return f.kick
}

func (f Force) GetHandling() float64 {
	return f.handling
}

//#region
func (f Force) AsGutCommandString() string {
	rotStr := fmt.Sprintf("%d", int(f.GetRot()))
	xStr := fmt.Sprintf("%d", int(f.GetX()))
	yStr := fmt.Sprintf("%d", int(f.GetY()))

	return fmt.Sprint("*", xStr, ",", yStr, ",", rotStr, ",", f.GetKick(), ",", f.GetHandling(), ",", f.reset, "#")
}
