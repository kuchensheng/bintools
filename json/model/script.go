package model

type Script struct {
	Script   string `json:"script"`
	Method   string `json:"method"`
	Language string `json:"language"`
	Args     []any  `json:"args"`
}
