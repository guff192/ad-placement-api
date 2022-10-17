package service

import placement "github.com/guff192/ad-placement-api"

type Service struct {
	Partners placement.PartnerArray
}

func NewService(partners placement.PartnerArray) *Service {
	return &Service{
		Partners: partners,
	}
}
