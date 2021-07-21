package lib

func MinMaxInt32(target [2]int32, x int32) [2]int32 {
	if target[0] > x {
		target[0] = x
	}
	if target[1] < x {
		target[1] = x
	}
	return target
}
