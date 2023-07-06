package util

import (
	"fmt"

	"github.com/pkoukk/tiktoken-go"
)

// EncodePromptBytePair - Encode string to byte pair
func EncodePromptBytePair(input []string, model string) []int {
	var buffer []int
	enc, _ := tiktoken.EncodingForModel(model)
	if enc != nil && input != nil {
		buffer = enc.Encode(fmt.Sprint(input), nil, nil)
	}
	return buffer
}
