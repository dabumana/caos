package util

import (
	"github.com/pkoukk/tiktoken-go"
)

// EncoderPromptToken - Encode string and calculate amount of tokens required
func EncodePromptToken(input []string, model string) int {
	enc, _ := tiktoken.EncodingForModel(model)
	total := enc.Encode(input[0], nil, nil)
	return len(total)
}
