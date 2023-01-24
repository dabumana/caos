// Package parameters section
package parameters

import "caos/util"

// Engine - Model parameter
var Engine = "text-davinci-003"

// Models - Available models for testing purposes
var Models []string

// Probabilities - Amount of probabilities designed for the request
var Probabilities int32 = util.ParseInt32("\u0031")

// Results - Amount of results designed for the request
var Results int32 = util.ParseInt32("\u0031")

// Temperature - Amount of temperature designed for the request
var Temperature float32 = util.ParseFloat32("\u0030\u002e\u0034")

// Topp - Amount of topp designed for the request
var Topp float32 = util.ParseFloat32("\u0031\u002e\u0030")

// Penalty - Amount of penalty threshold designed for the request
var Penalty float32 = util.ParseFloat32("\u0030\u002e\u0035")

// Frequency - Amount of penalty frequency threshold designed for the request
var Frequency float32 = util.ParseFloat32("\u0030\u002e\u0035")

// PromptCtx - Contextual prompt designed for the request
var PromptCtx []string

// MaxTokens - Amount of tokens assigned for the request
var MaxTokens int64 = util.ParseInt64("\u0032\u0030\u0034\u0038")

// Mode - Select between (Edit/Text/Code)
var Mode = "Text"

// IsLoading - Actually working in a request
var IsLoading = false

// IsConversational - Conversational mode
var IsConversational = false

// IsEditable - Editable completion mode
var IsEditable = false

// IsTraining - Fine-tunning training mode
var IsTraining = false

// IsNewSession - New session
var IsNewSession = true
