// Parameters section
package paramaters

import "caos/util"

// Global parameters
var Engine = "text-davinci-003"
var Probabilities int32 = util.ParseInt32("\u0031")
var Results int32 = util.ParseInt32("\u0031")
var Temperature float32 = util.ParseFloat32("\u0030\u002e\u0034")
var Topp float32 = util.ParseFloat32("\u0031\u002e\u0030")
var Penalty float32 = util.ParseFloat32("\u0030\u002e\u0035")
var Frequency float32 = util.ParseFloat32("\u0030\u002e\u0035")
var PromptCtx []string
var MaxTokens int64 = util.ParseInt64("\u0032\u0030\u0034\u0038")
var Mode = "Text"

// Modes
var IsLoading = false
var IsConversational = false
var IsEditable = true
var IsTraining = false

// New session
var IsNewSession = true
