package models_test

import (
	"fmt"
	"math"
	"testing"

	"harianugrah.com/brainfreeze/pkg/models"
)

func TestDegreeRadianConversion(t *testing.T) {
	tolerance := 0.01

	deg := models.Degree(180)
	if math.Abs(float64(deg.AsRadian())-3.14) > tolerance {
		t.Fatalf(fmt.Sprintln("Gagal menerjemahkan degree ke radian, expected 3.14159 got", deg.AsRadian()))
	}

	deg = models.Degree(360)
	if math.Abs(float64(deg.AsRadian())-6.28) > tolerance {
		t.Fatalf(fmt.Sprintln("Gagal menerjemahkan degree ke radian, expected 6.28319 got", deg.AsRadian()))
	}

	rad := models.Radian(1)
	if math.Abs(float64(rad.AsDegree())-57.2958) > tolerance {
		t.Fatalf(fmt.Sprintln("Gagal menerjemahkan radian ke degree, expected 57.2958 got", rad.AsDegree()))
	}
}

func TestDegreeHalfCircle(t *testing.T) {
	d := models.Degree(0)
	if float64(d) != 0 {
		t.Fatalf(fmt.Sprintln("Nilai yang diharapkan", 0, "nilai yang didapat", float64(d)))
	}

	d.Rotate(models.Degree(30))
	if float64(d) != 30 {
		t.Fatalf(fmt.Sprintln("Nilai yang diharapkan", 30, "nilai yang didapat", float64(d)))
	}

	d.Rotate(models.Degree(-45))
	if float64(d) != -15 {
		t.Fatalf(fmt.Sprintln("Nilai yang diharapkan", -15, "nilai yang didapat", float64(d)))
	}

	d.Rotate(models.Degree(360))
	if float64(d) != -15 {
		t.Fatalf(fmt.Sprintln("Nilai yang diharapkan", -15, "nilai yang didapat", float64(d)))
	}

	d.Rotate(models.Degree(-180))
	if float64(d) != 165 {
		t.Fatalf(fmt.Sprintln("Nilai yang diharapkan", 165, "nilai yang didapat", float64(d)))
	}

	d.Rotate(models.Degree(20))
	if float64(d) != -175 {
		t.Fatalf(fmt.Sprintln("Nilai yang diharapkan", -175, "nilai yang didapat", float64(d)))
	}

	d.Rotate(models.Degree(-20))
	if float64(d) != 165 {
		t.Fatalf(fmt.Sprintln("Nilai yang diharapkan", 165, "nilai yang didapat", float64(d)))
	}

	d.Rotate(models.Degree(-180))
	if float64(d) != -15 {
		t.Fatalf(fmt.Sprintln("Nilai yang diharapkan", -15, "nilai yang didapat", float64(d)))
	}

	d.Rotate(models.Degree(360 * 8))
	if float64(d) != -15 {
		t.Fatalf(fmt.Sprintln("Nilai yang diharapkan", -15, "nilai yang didapat", float64(d)))
	}

	d.Rotate(models.Degree(360 * -99))
	if float64(d) != -15 {
		t.Fatalf(fmt.Sprintln("Nilai yang diharapkan", -15, "nilai yang didapat", float64(d)))
	}
}
