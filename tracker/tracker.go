package tracker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type packageProgress struct {
	Status      string `json:"status"`
	StatusText  string `json:"statusText"`
	ImageStatus string `json:"imageStatus"`
	ImageText   string `json:"imageText"`
}

type glsResponse struct {
	TuStatus []struct {
		History []struct {
			Time    string `json:"time"`
			Date    string `json:"date"`
			Address struct {
				City        string `json:"city"`
				CountryCode string `json:"countryCode"`
				CountryName string `json:"countryName"`
			} `json:"address"`
			EvtDscr string `json:"evtDscr"`
		} `json:"history"`
		References []struct {
			Type  string `json:"type"`
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"references"`
		TuNo                   string `json:"tuNo"`
		ChangeDeliveryPossible bool   `json:"changeDeliveryPossible"`
		ProgressBar            struct {
			Level       int               `json:"level"`
			StatusText  string            `json:"statusText"`
			StatusBar   []packageProgress `json:"statusBar"`
			RetourFlag  bool              `json:"retourFlag"`
			EvtNos      []string          `json:"evtNos"`
			ColourIndex int               `json:"colourIndex"`
			StatusInfo  string            `json:"statusInfo"`
		} `json:"progressBar"`
		Owners []struct {
			Type string `json:"type"`
			Code string `json:"code"`
		} `json:"owners"`
		Infos []struct {
			Type  string `json:"type"`
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"infos"`
		DeliveryOwnerCode string `json:"deliveryOwnerCode"`
	} `json:"tuStatus"`
}

const glsURLPattern = "https://gls-group.eu/app/service/open/rest/HU/hu/rstt001?match=%s"

//const glsURLPattern = "https://gls-group.eu/app/service/open/rest/HU/en/rstt001?match=%s" // EN version

// ErrNotFound .
var ErrNotFound = errors.New("not found")

// Tracker .
type Tracker struct {
	client *http.Client
}

// New .
func New() *Tracker {
	return &Tracker{
		client: http.DefaultClient,
	}
}

type PackageStatus struct {
	StatusText string
	ImageText  string
}

func (t *Tracker) Get(packageID string) (PackageStatus, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(glsURLPattern, *pkg), nil)
	if err != nil {
		return PackageStatus{}, err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return PackageStatus{}, err
	}
	defer resp.Body.Close()

	var jsonResp glsResponse
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		return PackageStatus{}, err
	}

	if len(jsonResp.TuStatus) < 1 {
		return PackageStatus{}, ErrNotFound
	}

	status, ok := getCurrentStatus(jsonResp.TuStatus[0].ProgressBar.StatusBar)
	if !ok {
		return PackageStatus{}, errors.New("no current status found")
	} else {
		return PackageStatus{
			StatusText: status.StatusText,
			ImageText:  status.ImageText,
		}, nil
	}
}

func getCurrentStatus(pp []packageProgress) (packageProgress, bool) {
	for _, p := range pp {
		if p.ImageStatus == "CURRENT" {
			return p, true
		}
	}
	return packageProgress{}, false
}
