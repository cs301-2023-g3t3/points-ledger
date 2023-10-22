package models

type Input struct {
	Action string `json:"action"`
	Amount int    `json:"amount"`
}

type RequestMetadata struct {
    UserAgent string
    SourceIP  string
}