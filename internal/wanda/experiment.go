package wanda

import (
	"gocv.io/x/gocv"
)

// func gftt(gray *gocv.Mat) {
// 	mat := gocv.NewMat()
// 	gocv.GoodFeaturesToTrack(*gray, &mat, 20, 0.1, 10)
// 	// fmt.Println(mat.Rows(), mat.Cols(), mat.Channels(), mat.Total())
// 	fmt.Println(mat.GetUCharAt(0, 0))
// }

// func hl(mat *gocv.Mat) {
// 	// mat := gocv.IMRead("test-data/testdata20.jpg", gocv.IMReadGrayScale)
// 	matLines := gocv.NewMat()
// 	matCanny := gocv.NewMat()
// 	if mat.Empty() {
// 		fmt.Println("Did not load")
// 	}
// 	gocv.Canny(*mat, &matCanny, 100, 100)
// 	gocv.HoughLinesP(matCanny, &matLines, 0.5, math.Pi/360, 75)

// 	fmt.Println(matLines.Rows())

// 	for index1 := 0; index1 < matLines.Rows(); index1++ {
// 		pt1 := image.Pt(int(matLines.GetVeciAt(index1, 0)[0]), int(matLines.GetVeciAt(index1, 0)[1]))
// 		pt2 := image.Pt(int(matLines.GetVeciAt(index1, 0)[2]), int(matLines.GetVeciAt(index1, 0)[3]))
// 		gocv.Line(mat, pt1, pt2, color.RGBA{255, 255, 255, 1}, 10)
// 	}
// }

func whitte(mat *gocv.Mat, out *gocv.Mat) {
	upper := gocv.NewScalar(255, 255, 255, 1)
	lower := gocv.NewScalar(230, 230, 230, 0)

	gocv.InRangeWithScalar(*mat, lower, upper, out)

	// dilateMat := gocv.Ones(3, 3, gocv.MatTypeCV8UC1)
	// gocv.Dilate(*out, out, dilateMat)
	mat.CopyTo(out)
}
