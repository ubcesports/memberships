package util

import (
	"fmt"
	"math/rand/v2"
)

// GenerateNonStudentID produces a random student id in the format "N1234567",
// used for users who are not UBC students and therefore have no real student ID.
func GenerateNonStudentID() string {
	num := rand.IntN(10_000_000) // 0 to 9999999 (7 digits)
	return fmt.Sprintf("N%07d", num)
}
