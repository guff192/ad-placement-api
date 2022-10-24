package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	placement "github.com/guff192/ad-placement-api"
	"github.com/sirupsen/logrus"
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

func (s *ImpService) GetAdsForPlacements(id string, tiles []placement.Tile, context placement.Context) ([]placement.ImpResponse, error) {
	imps, err := s.getAllImps(id, tiles, context)
	if err != nil {
		return nil, err
	}

	mostExpensiveImps := s.findMostExpensiveImps(imps)

	var ads []placement.ImpResponse
	for _, tile := range tiles {
		id := tile.Id

		imp, ok := mostExpensiveImps[id]
		if ok {
			resp := &placement.ImpResponse{
				Id:     id,
				Width:  imp.Width,
				Height: imp.Height,
				Title:  imp.Title,
				URL:    imp.URL,
			}
			ads = append(ads, *resp)
		}
	}

	return ads, nil
}

func (s *ImpService) getAllImps(id string, tiles []placement.Tile, context placement.Context) ([]placement.Imp, error) {
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
	reqBodyReader := bytes.NewReader(reqBytes)
	url := partner.Addr + ":" + strconv.Itoa(partner.Port) + "/bid_request"
	request, err := http.NewRequest("POST", url, reqBodyReader)
	if err != nil {
		return
	}

	logrus.Info("Getting imps from partner: " + url)

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

func (s *ImpService) findMostExpensiveImps(imps []placement.Imp) map[uint]placement.Imp {
	// creating map for storing the most expensive imp for each id
	impMap := make(map[uint]placement.Imp)

	// filling in map with the most expensive imps
	for _, imp := range imps {
		id := imp.Id

		mostExpensiveImp, ok := impMap[id]
		if !ok {
			impMap[id] = imp
		}
		if imp.Price > mostExpensiveImp.Price {
			impMap[id] = imp
		}
	}

	return impMap
}
