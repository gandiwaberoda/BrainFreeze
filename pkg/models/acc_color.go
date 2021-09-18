package models

import (
	"image/color"

	"harianugrah.com/brainfreeze/internal/wanda/pxop"
)

type AcceptableColor struct {
	Id        int
	Name      string
	Upper     pxop.Vecb
	Lower     pxop.Vecb
	Visualize color.RGBA
}
