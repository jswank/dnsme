package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type apiDomainList struct {
	List []string `json:"list"`
}

type apiDomain struct {
	Name              string   `json:"name"`
	NameServers       []string `json:"nameServer"`
	VanityNameServers []string `json:"vanityNameServers"`
	GtdEnabled        bool     `json:"gtdEnabled"`
	Error             []string `json:"error,omitempty"`
}

type apiSecondary struct {
	Name  string   `json:"name"`
	IP    []string `json:"ip"`
	Error []string `json:"error,omitempty"`
}

type apiRecord struct {
	Name        string   `json:"name"`
	ID          int      `json:"id,omitemtpy"`
	Type        string   `json:"type"`
	Data        string   `json:"data"`
	GtdLocation string   `json:"gtdLocation"`
	TTL         int      `json:"ttl"`
	Password    string   `json:"password,omitempty"`
	Error       []string `json:"error,omitempty"`
}

func getDomainList() (domains apiDomainList, err error) {

	req, err := http.NewRequest("GET", api_url+"/domains/", nil)
	if err != nil {
		return
	}

	err = makeRequest(req, &domains)
	if err != nil {
		return
	}

	return
}

func getDomainInfo(domain string) (info apiDomain, err error) {

	req, err := http.NewRequest("GET", api_url+"/domains/"+domain, nil)
	if err != nil {
		return
	}

	err = makeRequest(req, &info)
	if err != nil {
		return
	}

	if len(info.Error) > 0 {
		errStr := strings.Join(info.Error, " ")
		err = errors.New(errStr)
		return
	}

	return

}

func addDomain(domain apiDomain) (domainResponse apiDomain, err error) {

	/* TODO: use json.Encoder or something providing an io.Reader? */
	jsonBody, err := json.Marshal(domain)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	buf.Write(jsonBody)
	// END TODO	

	req, err := http.NewRequest("PUT", api_url+"/domains/"+domain.Name, &buf)
	if err != nil {
		return
	}

	req.Header.Add("content-type", "application/json")

	err = makeRequest(req, &domainResponse)
	if err != nil {
		return
	}

	if len(domainResponse.Error) > 0 {
		errStr := strings.Join(domainResponse.Error, " ")
		err = errors.New(errStr)
		return
	}
	return
}

func deleteDomain(domain string) (err error) {

	req, err := http.NewRequest("DELETE", api_url+"/domains/"+domain, nil)
	if err != nil {
		return
	}

	err = makeRequest(req, nil)
	if err != nil {
		return
	}

	return

}

func getSecondaryList() (domains apiDomainList, err error) {

	req, err := http.NewRequest("GET", api_url+"/secondary/", nil)
	if err != nil {
		return
	}

	err = makeRequest(req, &domains)
	if err != nil {
		return
	}

	return
}

func getSecondary(domain string) (info apiSecondary, err error) {

	req, err := http.NewRequest("GET", api_url+"/secondary/"+domain, nil)
	if err != nil {
		return
	}

	err = makeRequest(req, &info)
	if err != nil {
		return
	}

	if len(info.Error) > 0 {
		errStr := strings.Join(info.Error, " ")
		err = errors.New(errStr)
		return
	}

	return

}

func addSecondary(s apiSecondary) (secondary apiSecondary, err error) {

	/* TODO: use json.Encoder or something providing an io.Reader? */
	jsonBody, err := json.Marshal(s)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	buf.Write(jsonBody)
	// END TODO	

	req, err := http.NewRequest("PUT", api_url+"/secondary/"+s.Name, &buf)
	if err != nil {
		return
	}

	req.Header.Add("content-type", "application/json")

	err = makeRequest(req, &secondary)
	if err != nil {
		return
	}

	if len(secondary.Error) > 0 {
		errStr := strings.Join(secondary.Error, " ")
		err = errors.New(errStr)
		return
	}
	return
}

func deleteSecondary(domain string) (err error) {

	req, err := http.NewRequest("DELETE", api_url+"/secondary/"+domain, nil)
	if err != nil {
		return
	}

	err = makeRequest(req, nil)
	if err != nil {
		return
	}

	return

}

func getDomainRecord(id, domain string) (record apiRecord, err error) {

	req, err := http.NewRequest("GET", api_url+"/domains/"+domain+"/records/"+id, nil)
	if err != nil {
		return
	}

	err = makeRequest(req, &record)
	if err != nil {
		return
	}

	if len(record.Error) > 0 {
		errStr := strings.Join(record.Error, " ")
		err = errors.New(errStr)
		return
	}

	// This is a shortcoming in the DNSME API: CNAME responses may have an 
	// empty "data" field, but updating/adding records always require the 
	// data field.
	if record.Type == "CNAME" && record.Data == "" {
		record.Data = domain + "."
	}

	return
}

func getDomainRecords(domain string, vals interface{}) (records []apiRecord, err error) {

	req, err := http.NewRequest("GET", api_url+"/domains/"+domain+"/records", nil)
	if err != nil {
		return
	}

	switch i := vals.(type) {
	case *url.Values:
		req.URL.RawQuery = i.Encode()
	}

	err = makeRequest(req, &records)
	if err != nil {
		return
	}

	for i, rec := range records {
		// This is a shortcoming in the DNSME API
		if rec.Type == "CNAME" && rec.Data == "" {
			records[i].Data = domain + "."
		}
	}
	return
}

func deleteDomainRecord(id, domain string) (err error) {

	req, err := http.NewRequest("DELETE", api_url+"/domains/"+domain+"/records/"+id, nil)
	if err != nil {
		return
	}

	err = makeRequest(req, nil)
	if err != nil {
		return
	}

	return

}

func addDomainRecord(domain string, r apiRecord) (record apiRecord, err error) {

	isUpdate := false

	/* TODO: use json.Encoder or something providing an io.Reader? */
	jsonBody, err := json.Marshal(r)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	buf.Write(jsonBody)
	// END TODO	

	// whether this is an "add" or "update" depends on the value of the "ID" field
	var method, url string
	if r.ID == 0 { // add
		method = "POST"
		url = api_url + "/domains/" + domain + "/records/"
	} else { // update
		method = "PUT"
		url = api_url + "/domains/" + domain + "/records/" + strconv.Itoa(r.ID)
		isUpdate = true
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return
	}

	req.Header.Add("content-type", "application/json")

	err = makeRequest(req, &record)
	if err != nil {
		return
	}

	// update requests return an empty body
	if isUpdate {
		return
	}

	if len(record.Error) > 0 {
		errStr := strings.Join(record.Error, " ")
		err = errors.New(errStr)
		return
	}

	return
}

func addDnsmeHeaders(r *http.Request) {
	r.Header.Add("x-dnsme-apiKey", api_key)

	requestDate := time.Now().UTC().Format(time.RFC1123)
	r.Header.Add("x-dnsme-requestDate", requestDate)

	h := hmac.New(sha1.New, []byte(secret_key))
	h.Write([]byte(requestDate))
	r.Header.Add("x-dnsme-hmac", fmt.Sprintf("%x", h.Sum(nil)))

	r.Header.Add("Accept", "application/json")

	return
}

/*
 * makeRequest() performs http requests that are built by API functions.
 * It updates the global requestsRemaining based on the API response, and 
 * uses a simple retry mechanism whenever the API rate limit has been 
 * exceeded.
 */
func makeRequest(r *http.Request, into interface{}) (err error) {
	var resp *http.Response

	// make this part of the struct?
	client := &http.Client{}

	addDnsmeHeaders(r)

	max_tries := 10

	for t := 0; t < max_tries; t++ {
		if debug {
			dump, d_err := httputil.DumpRequestOut(r, true)
			if d_err == nil {
				os.Stderr.Write(dump)
			}
		}
		resp, err = client.Do(r)
		requestsRemaining, _ = strconv.Atoi(resp.Header.Get("x-dnsme-requestsRemaining"))
		if err != nil && requestsRemaining > 0 {
			// fmt.Printf("Error performing HTTP request, %s", err)
			return
		}
		if requestsRemaining == 0 {
			//fmt.Fprintf(os.Stderr, "API rate-limit exceeded, sleeping for 20 seconds (try %d of %d)\n", t, max_tries)
			time.Sleep(30 * time.Second) // 6 seconds
		} else {
			break
		}
	}
	if debug {
		dump, d_err := httputil.DumpResponse(resp, true)
		if d_err == nil {
			os.Stderr.Write(dump)
		}
	}
	if resp.StatusCode == http.StatusForbidden {
		err = errors.New("API access forbidden")
		return
	}
	if resp.StatusCode == http.StatusNotFound {
		err = errors.New("Not found")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, into)
	if err != nil {
		return
	}

	return
}
