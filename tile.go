package placement

import (
	"math"
)

type Tile struct {
	Id    uint    `json:"id"`
	Width uint    `json:"width"`
	Ratio float64 `json:"ratio"`
}

func (t *Tile) ToImpRequest() *ImpRequest {
	id := t.Id
	minwidth := t.Width
	minheight := math.Floor(float64(t.Width) * t.Ratio)

	return &ImpRequest{
		Id:        id,
		Minwidth:  minwidth,
		Minheight: uint(minheight),
	}
}
