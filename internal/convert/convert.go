package convert

import (
	"advertisement/model"
	"math"
)

func ConvertPartnerToAnswer(impPart model.ImpPart) model.ImpAns {
	return model.ImpAns{
		ID:     impPart.ID,
		Width:  impPart.Width,
		Height: impPart.Height,
		Title:  impPart.Title,
		Url:    impPart.Url,
	}
}

func ConvertTilesToImp(tiles model.Tiles) model.Imp {
	return model.Imp{
		ID:        tiles.ID,
		MinWidth:  tiles.Width,
		MinHeight: uint(math.Floor(float64(tiles.Width) * float64(tiles.Ratio))),
	}
}
