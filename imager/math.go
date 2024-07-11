package imager

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// lcm is used to calculate the total number of frames in an animation,
// when there are multiple animation layers with different frame counts.
func lcm(a, b int) int {
	return a * b / gcd(a, b)
}
