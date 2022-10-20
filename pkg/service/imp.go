package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	placement "github.com/guff192/ad-placement-api"
)

type ImpService struct {
	Partners placement.PartnerArray
}

func NewImpService(partners placement.PartnerArray) *ImpService {
	return &ImpService{
		Partners: partners,
	}
}

type BidRequest struct {
	Id      string                 `json:"id"`
	Imp     []placement.ImpRequest `json:"imp"`
	Context placement.Context      `json:"context"`
}

func (s *ImpService) GetAllImps(id string, tiles []placement.Tile, context placement.Context) ([]placement.Imp, error) {
	var reqImps []placement.ImpRequest
	for _, tile := range tiles {
		imp := tile.ToImpRequest()
		reqImps = append(reqImps, *imp)
	}

	// creating request to partners
	request := &BidRequest{
		Id:      id,
		Imp:     reqImps,
		Context: context,
	}
	reqBytes, err := json.Marshal(*request)
	if err != nil {
		return nil, err
	}

	// collecting imps from partners
	var imps []placement.Imp
	client := &http.Client{
		Timeout: 250 * time.Millisecond,
	}
	for _, partner := range s.Partners {
		go s.getImpsFromAddr(client, partner, reqBytes, &imps)
	}

	return imps, nil
}

type impPartnerResponse struct {
	Id  string          `json:"id"`
	Imp []placement.Imp `json:"imp"`
}

func (s *ImpService) getImpsFromAddr(client *http.Client, partner placement.PartnerAddr, reqBytes []byte, imps *[]placement.Imp) {
	// creating request
	addr, port, reqBodyReader := partner.Addr, partner.Port, bytes.NewReader(reqBytes)
	request, err := http.NewRequest("POST", addr+":"+string(port), reqBodyReader)
	if err != nil {
		return
	}

	// getting the response and checking Content-Type
	response, err := client.Do(request)
	if err != nil || response.Header.Get("Content-Type") != "application/json" {
		return
	}

	// reading body bytes and unmarshalling to var
	var impResponse impPartnerResponse
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(bodyBytes, &impResponse); err != nil {
		return
	}

	*imps = append(*imps, impResponse.Imp...)
}
