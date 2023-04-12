// Package model section
package model

// PredictProperties - Request model
type PredictProperties struct {
	Input   []string        `json:"context"`
	Details PredictResponse `json:"details"`
}

// PredictResponse - Response model
type PredictResponse struct {
	Documents []Predict `json:"documents"`
}

// Predict - Document properties model
type Predict struct {
	AverageProb       float64     `json:"average_generated_prob"`
	CompletelyProb    float64     `json:"completely_generated_prob"`
	OverallBurstiness float64     `json:"overall_burstiness"`
	Sentences         []Sentence  `json:"sentences"`
	Paragraphs        []Paragraph `json:"paragraphs"`
}

// Sentence - Nested sentence model
type Sentence struct {
	Sentence      string  `json:"sentence"`
	Perplexity    float64 `json:"perplexity"`
	GeneratedProb float64 `json:"generated_prob"`
}

// Paragraph - Nested paragraphs model
type Paragraph struct {
	Index           int     `json:"start_sentence_index"`
	NumberSentences int     `json:"num_sentences"`
	CompletelyProb  float64 `json:"completely_generated_prob"`
}

// PredictRequest - Request model
type PredictRequest struct {
	Document string `json:"document"`
}
