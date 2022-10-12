package model

type Answer struct {
	ID     string   `json:"id"`
	ImpAns []ImpAns `json:"imp"`
}

type ImpAns struct {
	ID     uint    `json:"id"`
	Width  uint    `json:"width"`
	Height float32 `json:"height"`
	Title  string  `json:"title"`
	Url    string  `json:"url"`
}
