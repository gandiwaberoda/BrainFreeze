package models_test

import (
	"fmt"
	"testing"

	"harianugrah.com/brainfreeze/pkg/models"
)

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
