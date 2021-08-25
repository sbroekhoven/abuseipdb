package abuseipdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

	values := url.Values{}
	values.Set("ip", ipAddress)
	values.Set("categories", categories)
	values.Set("comment", comment)

	request, err := http.NewRequest("POST", APIEndpoint, bytes.NewBufferString(values.Encode()))
	if err != nil {
		return ReportResponse{}, err
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/token", siteHost), bytes.NewBufferString(values.Encode()))
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value") // This makes it work

	request.Header.Add("Key", c.APIKey)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")
	request.Header.Add("User-Agent", "Check AbuseIPDB by github.com/binaryfigments")

	resp, err := client.Do(request)
	if err != nil {
		return ReportResponse{}, err
	}

	fmt.Println(resp.Status)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ReportResponse{}, err
	}

	// var response ReportResponse
	response := ReportResponse{}
	json.Unmarshal(body, &response)

	// For testing
	fmt.Println(string(body))

	return response, err
}
