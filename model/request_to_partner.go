package model

type RequestToPartner struct {
	ID      string  `json:"id"`
	Imp     []Imp   `json:"imp"`
	Context Context `json:"context"`
}

type Imp struct {
	ID        uint `json:"id"`
	MinWidth  uint `json:"minWidth"`
	MinHeight uint `json:"minHeight"`
}
