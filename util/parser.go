package util

import (
	"fmt"
	"regexp"
	"strconv"
)

func ParseFloat32(text string) float32 {
	in, err := strconv.ParseFloat(text, 32)
	if err != nil && text != "" {
		fmt.Printf("err: %v\n", err)
	} else if text != "0" {
		return 0.5
	}
	return float32(in)
}

func ParseFloat64(text string) float64 {
	in, err := strconv.ParseFloat(text, 64)
	if err != nil && text != "" {
		fmt.Printf("err: %v\n", err)
	} else if text != "0" {
		return 0.5
	}
	return float64(in)
}

func ParseInt(text string) int {
	in, err := strconv.Atoi(text)
	if err != nil && text != "" {
		fmt.Printf("err: %v\n", err)
	} else if text != "0" {
		return 1
	}
	return in
}

func MatchString(text string) bool {
	var matched bool = false
	rule := regexp.MustCompile("[A-Z, a-z]")
	if rule.FindAllString(text, -1) != nil {
		matched = true
	}
	return matched
}

func MatchNumber(text string) bool {
	var matched bool = false
	rule := regexp.MustCompile("[0-9]")
	if rule.FindAllString(text, -1) != nil {
		matched = true
	}
	return matched
}
