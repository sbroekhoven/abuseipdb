package abuseipdb

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Documentation: https://docs.abuseipdb.com/#check-endpoint

type CheckResponse struct {
	Data checkData `json:"data"`
}

type checkData struct {
	IpAddress            string        `json:"ipAddress"`
	IsPublic             bool          `json:"isPublic"`
	IpVersion            int           `json:"ipVersion"`
	IsWhitelisted        bool          `json:"isWhitelisted"`
	AbuseConfidenceScore int           `json:"abuseConfidenceScore"`
	CountryCode          string        `json:"countryCode"`
	CountryName          string        `json:"countryName"`
	UsageType            string        `json:"usageType"`
	Isp                  string        `json:"isp"`
	Domain               string        `json:"domain"`
	Hostnames            []string      `json:"hostnames"`
	TotalReports         int           `json:"totalReports"`
	NumDistinctUsers     int           `json:"numDistinctUsers"`
	LastReportedAt       string        `json:"lastReportedAt"`
	Reports              []checkReport `json:"reports"`
}

type checkReport struct {
	ReportedAt          string `json:"reportedAt"`
	Comment             string `json:"comment"`
	Categories          []int  `json:"categories"`
	ReporterId          int    `json:"reporterId"`
	ReporterCountryCode string `json:"reporterCountryCode"`
	ReporterCountryName string `json:"reporterCountryName"`
}

func Check(c *Configuration, ipAddress string, maxAgeInDays int, verbose bool) (CheckResponse, error) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", "https://api.abuseipdb.com/api/v2/check", nil)
	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Add("Key", c.APIKey)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "Check AbuseIPDB by github.com/binaryfigments")

	query := request.URL.Query()
	query.Add("ipAddress", ipAddress)
	query.Add("maxAgeInDays", strconv.Itoa(maxAgeInDays))
	query.Add("verbose", "")
	request.URL.RawQuery = query.Encode()

	resp, err := client.Do(request)
	if err != nil {
		return CheckResponse{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return CheckResponse{}, err
	}

	// var response CheckResponse
	response := CheckResponse{}
	json.Unmarshal(body, &response)

	return response, err
}
