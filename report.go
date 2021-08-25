package abuseipdb

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Documentation: https://docs.abuseipdb.com/#report-endpoint

type ReportResponse struct {
	Data  reportData  `json:"data,omitempty"`
	Error reportError `json:"error,omitempty"`
}

type reportData struct {
	IpAddress            string `json:"ipAddress,omitempty"`
	AbuseConfidenceScore int    `json:"abuseConfidenceScore,omitempty"`
}

type reportError struct {
	Detail string            `json:"detail,omitempty"`
	Status int               `json:"status,omitempty"`
	Source reportErrorSource `json:"source,omitempty"`
}

type reportErrorSource struct {
	Parameter string `json:"parameter,omitempty"`
}

// Report function to report 1 IP address to AbuseIPDB
func Report(c *Configuration, ipAddress string, categories string, comment string) (ReportResponse, error) {
	APIEndpoint := c.APIURL + "/report"

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("POST", APIEndpoint, nil)
	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Add("Key", c.APIKey)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "Check AbuseIPDB by github.com/binaryfigments")

	query := request.URL.Query()
	query.Add("ipAddress", ipAddress)
	query.Add("categories", categories)
	query.Add("comment", comment)
	request.URL.RawQuery = query.Encode()

	resp, err := client.Do(request)
	if err != nil {
		return ReportResponse{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ReportResponse{}, err
	}

	// var response ReportResponse
	response := ReportResponse{}
	json.Unmarshal(body, &response)

	return response, err
}
