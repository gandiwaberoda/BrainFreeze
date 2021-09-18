package gutmodel

import (
	"errors"
	"strconv"
	"strings"

	"harianugrah.com/brainfreeze/pkg/models"
)

type Araya struct {
	Dist0 models.Centimeter // 0 derajat aka depan robot
}

func ParseAraya(str string) (Araya, error) {
	if str[0] != '*' || str[len(str)-1] != '#' {
		return Araya{}, errors.New("wrong araya format")
	}

	cleanHeaderTail := strings.ReplaceAll(str, "*", "")
	cleanHeaderTail = strings.ReplaceAll(cleanHeaderTail, "#", "")
	splitted := strings.Split(cleanHeaderTail, ",")

	if len(splitted) != 1 {
		return Araya{}, errors.New("wrong number of araya data")
	}

	_dist0, errDist0 := strconv.ParseFloat(splitted[0], 64)
	if errDist0 != nil {
		return Araya{}, errDist0
	}

	return Araya{
		Dist0: models.Centimeter(_dist0),
	}, nil
}
