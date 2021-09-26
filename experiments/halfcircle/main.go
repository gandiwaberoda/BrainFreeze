package main

import (
	"fmt"
)

func AsHalfCircle(r float64) int {
	y := int(r) % 360
	if y >= 0 && y <= 180 {
		return (y)
	} else if y > 180 && y < 360 {
		return (-180 + (y - 180))
	} else if y <= 0 && y >= -180 {
		return (y)
	} else if y < -180 {
		return 180 - ((-1 * y) - 180)
	}

	return int(r)
}

func main() {
	// fmt.Println("0", AsHalfCircle(0))
	// fmt.Println("100", AsHalfCircle(100))
	// fmt.Println("-100", AsHalfCircle(-100))
	// fmt.Println("-190", AsHalfCircle(-190))
	// fmt.Println("190", AsHalfCircle(190))
	// fmt.Println("190", AsHalfCircle(190))

	for i := -720; i <= 720; i += 10 {
		fmt.Println(i, AsHalfCircle(float64(i)))
	}
}
