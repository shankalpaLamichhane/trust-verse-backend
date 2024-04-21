package dto

type SentimentResponse struct {
	RobertaNeg float64 `json:"roberta_neg"`
	RobertaPos float64 `json:"roberta_pos"`
	RobertaNeu float64 `json:"roberta_neu"`
}

type FakeNewsResponse struct {
	Prediction string `json:"prediction"`
}

type HealthResponse struct {
	Prediction      string  `json:"prediction"`
	PredictionScore float64 `json:"prediction_score"`
}
