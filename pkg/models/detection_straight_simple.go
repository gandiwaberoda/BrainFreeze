package models

type StraightDetection struct {
	Length        int
	UpperY        int
	LowerY        int
	DetectedColor AcceptableColor
}
type StraightDetectionObj struct {
	ClosestDistPx     int
	FurthestDistPx    int
	DetectedColorName string
}
