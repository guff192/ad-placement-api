package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"
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

	logrus.Info("Collected imps from partners for request [" + id + "]")
	logrus.Debug(imps)

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

	logrus.Debug(ads)

	return ads, nil
}

func (s *ImpService) getAllImps(id string, tiles []placement.Tile, context placement.Context) ([]placement.Imp, error) {
	var reqImps []placement.ImpRequest
	for _, tile := range tiles {
		imp := tile.ToImpRequest()
		reqImps = append(reqImps, *imp)
	}

	logrus.Debug(reqImps)

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
	impCh := make(chan []placement.Imp, len(s.Partners))
	var wg sync.WaitGroup // wait group for waiting all responses to be done
	client := &http.Client{
		Timeout: 250 * time.Millisecond,
	}
	for _, partner := range s.Partners {
		wg.Add(1)
		go s.getImpsFromAddr(client, partner, reqBytes, impCh, &wg)
	}
	wg.Wait() // waiting for all requests to be done

	close(impCh)

	// collecting all results into one slice
	var impResult []placement.Imp
	for imps := range impCh {
		logrus.Info("Parsing partner response")
		impResult = append(impResult, imps...)
	}

	return impResult, nil
}

type impPartnerResponse struct {
	Id  string          `json:"id"`
	Imp []placement.Imp `json:"imp"`
}

func (s *ImpService) getImpsFromAddr(client *http.Client, partner placement.PartnerAddr, reqBytes []byte, imps chan []placement.Imp, wg *sync.WaitGroup) {
	// decrement waitgroup counter when done
	defer wg.Done()

	// creating request
	reqBodyReader := bytes.NewReader(reqBytes)
	url := "http://" + partner.Addr + ":" + strconv.Itoa(partner.Port) + "/bid_request"
	request, err := http.NewRequest("POST", url, reqBodyReader)
	if err != nil {
		logrus.Warn("Error while creating request: " + err.Error())
		return
	}

	logrus.Info("Getting imps from partner: " + url)

	// getting the response and checking Content-Type
	response, err := client.Do(request)
	if err != nil {
		logrus.Warn("Error while getting partner response: " + err.Error())
		return
	}

	// reading body bytes and unmarshalling to var
	var impResponse impPartnerResponse
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Warn("Error while reading response body: " + err.Error())
		return
	}
	if err = json.Unmarshal(bodyBytes, &impResponse); err != nil {
		logrus.Warn("Error while unmarshalling response body: " + err.Error())
		return
	}

	logrus.Info("\nGot response from partner <" + url + ">")
	logrus.Debug(impResponse)

	imps <- impResponse.Imp
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
