package utils

func NewFloat32(f float32) *float32 {
	return &f
}

func CalcLevel(experience int64) int32 {
	return int32(experience / 100)
}
