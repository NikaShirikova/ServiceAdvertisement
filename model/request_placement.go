package model

type RequestPlacement struct {
	ID      string  `json:"id" validate:"required"`
	Tiles   []Tiles `json:"tiles" validate:"required"`
	Context Context `json:"context" validate:"required"`
}

type Tiles struct {
	ID    uint    `json:"id" validate:"required"`
	Width uint    `json:"width" validate:"required"`
	Ratio float32 `json:"ratio" validate:"required"`
}

type Context struct {
	IP        string `json:"ip" validate:"required"`
	UserAgent string `json:"userAgent" validate:"required"`
}
