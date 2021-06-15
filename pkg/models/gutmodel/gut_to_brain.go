package gutmodel

import (
	"errors"
	"strconv"
	"strings"

	"harianugrah.com/brainfreeze/pkg/models"
)

type GutToBrain struct {
	AbsX        models.Centimeter
	AbsY        models.Centimeter
	Gyro        models.Degree
	IsDribbling bool
}

func ParseGutToBrain(str string) (GutToBrain, error) {
	if str[0] != '*' || str[len(str)-1] != '#' {
		return GutToBrain{}, errors.New("wrong format")
	}

	cleanHeaderTail := strings.ReplaceAll(str, "*", "")
	cleanHeaderTail = strings.ReplaceAll(cleanHeaderTail, "#", "")
	splitted := strings.Split(cleanHeaderTail, ",")

	if len(splitted) != 4 {
		return GutToBrain{}, errors.New("wrong number of data")
	}

	absX, errAbsX := strconv.ParseFloat(splitted[0], 64)
	if errAbsX != nil {
		return GutToBrain{}, errAbsX
	}

	absY, errAbsY := strconv.ParseFloat(splitted[1], 64)
	if errAbsY != nil {
		return GutToBrain{}, errAbsY
	}

	gyro, errGyro := strconv.ParseFloat(splitted[2], 64)
	if errGyro != nil {
		return GutToBrain{}, errGyro
	}

	isDribbling, errIsDribbling := strconv.ParseFloat(splitted[3], 64)
	if errIsDribbling != nil {
		return GutToBrain{}, errIsDribbling
	}

	isDribblingBool := isDribbling == 1.0

	return GutToBrain{
		AbsX:        models.Centimeter(absX),
		AbsY:        models.Centimeter(absY),
		Gyro:        models.Degree(gyro),
		IsDribbling: isDribblingBool,
	}, nil
}
