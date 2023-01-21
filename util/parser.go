// Internal utilities section
package util

import (
	"fmt"
	"regexp"
	"strconv"
)

// Parse string to float32
func ParseFloat32(text string) float32 {
	in, err := strconv.ParseFloat(text, 32)
	if err != nil && text != "" {
		fmt.Printf("err: %v\n", err)
		return 0.5
	}
	return float32(in)
}

// Parse string float64
func ParseFloat64(text string) float64 {
	in, err := strconv.ParseFloat(text, 64)
	if err != nil && text != "" {
		fmt.Printf("err: %v\n", err)
		return 0.5
	}
	return float64(in)
}

// Parse string int64
func ParseInt64(text string) int64 {
	in, err := strconv.ParseInt(text, 0, 64)
	if err != nil && text != "" {
		fmt.Printf("err: %v\n", err)
		return 1
	}
	return int64(in)
}

// Parse string int32
func ParseInt32(text string) int32 {
	in, err := strconv.ParseInt(text, 0, 32)
	if err != nil && text != "" {
		fmt.Printf("err: %v\n", err)
		return 1
	}
	return int32(in)
}

// Match string with regex compatibility (only letters from a-Z)
func MatchString(text string) bool {
	var matched bool = false
	rule := regexp.MustCompile("[A-Z, a-z]")
	if rule.FindAllString(text, -1) != nil {
		matched = true
	}
	return matched
}

// Match number with regex compatibility (only numbers)
func MatchNumber(text string) bool {
	var matched bool = false
	rule := regexp.MustCompile("[0-9]")
	if rule.FindAllString(text, -1) != nil {
		matched = true
	}
	return matched
}
