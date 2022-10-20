package service

import placement "github.com/guff192/ad-placement-api"

type Imp interface {
	GetAdsForPlacements(id string, tiles []placement.Tile, context placement.Context) ([]placement.ImpResponse, error)
}

type Service struct {
	Imp
}

func NewService(partners placement.PartnerArray) *Service {
	return &Service{
		Imp: NewImpService(partners),
	}
}
