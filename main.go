package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
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

func main() {
	pkg := flag.String("pkg", "", "Package ID")
	flag.Parse()

	if *pkg == "" {
		log.Fatalln("Set a package id please!")
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(glsURLPattern, *pkg), nil)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var jsonResp glsResponse
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		log.Fatalln(err)
	}

	if len(jsonResp.TuStatus) < 1 {
		log.Fatalln("No package found")
	}

	status, ok := getCurrentStatus(jsonResp.TuStatus[0].ProgressBar.StatusBar)
	if !ok {
		fmt.Println("No current status found")
	} else {
		fmt.Printf("Now: %s (%s)\n", status.StatusText, status.ImageText)
	}

	fmt.Println("\nHistory:")
	for _, s := range jsonResp.TuStatus[0].History {
		fmt.Printf("%s %s %s (%s %s)\n", s.EvtDscr, s.Date, s.Time, s.Address.CountryName, s.Address.City)
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
