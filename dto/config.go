package dto

type ConfigDto struct {
	Id        int    `json:"id"`
	Timeout   int64  `json:"timeout"`
	LabelName string `json:"label_name"`
}
