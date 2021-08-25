package abuseipdb

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Documentation: https://docs.abuseipdb.com/#blacklist-endpoint

type BlacklistResponse struct {
	Meta blacklistMeta   `json:"meta"`
	Data []blacklistData `json:"data"`
}

type BlacklistPlaintextResponse struct {
	XGeneratedAt string   `json:"x-generated-at"`
	IpAddresses  []string `json:"ip"`
}

type blacklistMeta struct {
	GeneratedAt string `json:"generatedAt"`
}

type blacklistData struct {
	IpAddress            string `json:"ipAddress"`
	AbuseConfidenceScore int    `json:"abuseConfidenceScore"`
	LastReportedAt       string `json:"lastReportedAt"`
}

// Get IP's to blacklist from abuseIPDB.com
func Blacklist(c *Configuration, confidenceMinimum int) (BlacklistResponse, error) {
	APIEndpoint := c.APIURL + "/blacklist"

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", APIEndpoint, nil)
	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Add("Key", c.APIKey)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "Check AbuseIPDB by github.com/binaryfigments")

	query := request.URL.Query()
	query.Add("confidenceMinimum", strconv.Itoa(confidenceMinimum))
	request.URL.RawQuery = query.Encode()

	resp, err := client.Do(request)
	if err != nil {
		return BlacklistResponse{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return BlacklistResponse{}, err
	}

	// var response BlacklistResponse
	response := BlacklistResponse{}
	json.Unmarshal(body, &response)

	return response, err
}

// Get IP's to blacklist from abuseIPDB.com
// TODO: Maybe real text output and not JSON.
func BlacklistPlaintext(c *Configuration, confidenceMinimum int, limit int) (BlacklistPlaintextResponse, error) {
	APIEndpoint := c.APIURL + "/blacklist"

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", APIEndpoint, nil)
	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Add("Key", c.APIKey)
	request.Header.Add("Accept", "text/plain")
	request.Header.Add("User-Agent", "Check AbuseIPDB by github.com/binaryfigments")

	query := request.URL.Query()
	query.Add("confidenceMinimum", strconv.Itoa(confidenceMinimum))
	if limit > 0 {
		query.Add("limit", strconv.Itoa(limit))
	}
	request.URL.RawQuery = query.Encode()

	resp, err := client.Do(request)
	if err != nil {
		return BlacklistPlaintextResponse{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return BlacklistPlaintextResponse{}, err
	}
	bodyips := string(body)

	var ips []string
	scanner := bufio.NewScanner(strings.NewReader(bodyips))
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		ips = append(ips, scanner.Text())
	}

	// var response BlacklistResponse
	response := BlacklistPlaintextResponse{}
	response.IpAddresses = ips
	response.XGeneratedAt = resp.Header.Get("X-Generated-At")
	json.Unmarshal(body, &response)

	return response, err
}

// X-Generated-At
// fmt.Println(resp.Header.Get("X-Generated-At"))
