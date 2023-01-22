// Package util section
package util

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

// ParseFloat32 - Parse string to float32
func ParseFloat32(text string) float32 {
	in, err := strconv.ParseFloat(text, 32)
	if err != nil && text != "" {
		fmt.Printf("err: %v\n", err)
		return 0.5
	}
	return float32(in)
}

// ParseFloat64 - Parse string float64
func ParseFloat64(text string) float64 {
	in, err := strconv.ParseFloat(text, 64)
	if err != nil && text != "" {
		fmt.Printf("err: %v\n", err)
		return 0.5
	}
	return float64(in)
}

// ParseInt64 - Parse string int64
func ParseInt64(text string) int64 {
	in, err := strconv.ParseInt(text, 0, 64)
	if err != nil && text != "" {
		fmt.Printf("err: %v\n", err)
		return 1
	}
	return int64(in)
}

// ParseInt32 - Parse string int32
func ParseInt32(text string) int32 {
	in, err := strconv.ParseInt(text, 0, 32)
	if err != nil && text != "" {
		fmt.Printf("err: %v\n", err)
		return 1
	}
	return int32(in)
}

// MatchString - Match string with regex compatibility (only letters from a-Z)
func MatchString(text string) bool {
	var matched = false
	rule := regexp.MustCompile("[A-Z, a-z]")
	if rule.FindAllString(text, -1) != nil {
		matched = true
	}
	return matched
}

// MatchNumber - Match number with regex compatibility (only numbers)
func MatchNumber(text string) bool {
	var matched = false
	rule := regexp.MustCompile("[0-9]")
	if rule.FindAllString(text, -1) != nil {
		matched = true
	}
	return matched
}

// ConstructPathFileToJSON - Initialize a directory for further storage in a JSON file
func ConstructPathFileToJSON(path string) *os.File {
	var dir string
	if dir, e := os.Getwd(); e != nil {
		fmt.Printf("dir: %v\n", dir)
	}

	if _, err := os.Stat(fmt.Sprintf("%s/%s", dir, path)); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	var out *os.File
	now := fmt.Sprint(time.Now().UTC())
	tsFile := fmt.Sprintf("%s-%s.json", path, now)
	pathOutput := filepath.Join(dir, path, tsFile)

	if _, err := os.Stat(fmt.Sprintf("%s/%s/%s", dir, path, tsFile)); os.IsNotExist(err) {
		out, _ = os.Create(pathOutput)
	} else {
		out, _ = os.OpenFile(pathOutput, 0, 0644)
	}

	return out
}

// ConstructPathFileToTXT - Initialize a directory for further storage in a TXT file
func ConstructPathFileToTXT(path string) *os.File {
	var dir string
	if dir, e := os.Getwd(); e != nil {
		fmt.Printf("dir: %v\n", dir)
	}

	if _, err := os.Stat(fmt.Sprintf("%s/%s", dir, path)); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	var out *os.File
	now := fmt.Sprint(time.Now().UTC())
	tsFile := fmt.Sprintf("%s-%s.txt", path, now)
	pathOutput := filepath.Join(dir, path, tsFile)

	if _, err := os.Stat(fmt.Sprintf("%s/%s/%s", dir, path, tsFile)); os.IsNotExist(err) {
		out, _ = os.Create(pathOutput)
	} else {
		out, _ = os.OpenFile(pathOutput, 0, 0644)
	}

	return out
}
