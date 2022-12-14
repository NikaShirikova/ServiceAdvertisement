package model

type PartnerRequest struct {
	ID      string    `json:"id"`
	ImpPart []ImpPart `json:"imp"`
}

type ImpPart struct {
	ID     uint    `json:"id"`
	Width  uint    `json:"width"`
	Height float32 `json:"height"`
	Title  string  `json:"title"`
	Url    string  `json:"url"`
	Price  float32 `json:"price"`
}
