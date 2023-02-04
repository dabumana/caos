// Package model section
package model

// PredictProperties - Request model
type PredictProperties struct {
	Input   []string        `json:"input"`
	Details PredictResponse `json:"details"`
}

// PredictRequest - Request model
type PredictRequest struct {
	Document string `json:"document"`
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
	Paragraphs        []Paragraph `json:"paragraphs"`
	Sentences         []Sentence  `json:"sentences"`
}

// Sentence - Nested sentence model
type Sentence struct {
	GeneratedProb int    `json:"generated_prob"`
	Perplexity    int    `json:"perplexity"`
	Sentence      string `json:"sentence"`
}

// Paragraph - Nested paragraphs model
type Paragraph struct {
	CompletelyProb  float64 `json:"completely_generated_prob"`
	NumberSentences int     `json:"num_sentences"`
	Index           int     `json:"start_sentence_index"`
}
