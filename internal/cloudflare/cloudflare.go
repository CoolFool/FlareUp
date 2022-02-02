package cloudflare

import (
	"bytes"
	"encoding/json"
	"errors"
	"flareup/internal/logging"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

var log = logging.Init()

const cfApiEndpoint = "https://api.cloudflare.com/client/v4/zones/"

type dnsRecord struct {
	Result []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	}
	Errors []struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
	Messages []struct {
	} `json:"messages"`
	Success bool `json:"success"`
}

type updateResponse struct {
	Result struct {
		Id      string `json:"id"`
		Content string `json:"content"`
	}
	Errors []struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
	Messages []interface {
	} `json:"messages"`
	Success bool `json:"success"`
}

type newRecord struct {
	Kind    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
	Ttl     int    `json:"ttl"`
}

type cloudflare struct {
	apiToken   string
	hostname   string
	domain     string
	zoneId     string
	recordType string
	recordId   string
	content    string
	proxied    bool
}

func requestMaker(client *http.Client, header http.Header, endpoint string) (body []byte) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	req.Header = header
	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
			return
		}
	}(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)

	}
	return
}

func (cf *cloudflare) setZoneId(client *http.Client, header http.Header) (err error) {
	base, err := url.Parse(cfApiEndpoint)
	if err != nil {
		log.Error(err)
		return errors.New("error while parsing cloudflare endpoint")
	}
	endpoint := base.String()
	params := url.Values{}
	params.Add("name", cf.domain)
	params.Add("status", "active")
	base.RawQuery = params.Encode()
	var data = new(dnsRecord)
	resp := requestMaker(client, header, endpoint)
	err = json.Unmarshal(resp, &data)
	if err != nil {
		log.Error(err)
		return
	}
	if len(data.Errors) > 0 || (len(data.Result) <= 0 && len(data.Errors) >= 0) {
		for _, err := range data.Errors {
			log.Error(fmt.Sprintf("%d %s while updating domain: %s", err.Code, err.Message, cf.domain))
		}
		return errors.New(fmt.Sprintf("Zone id for %s not found", cf.domain))
	}
	cf.zoneId = data.Result[0].Id
	return nil
}

func (cf *cloudflare) setRecord(client *http.Client, header http.Header) (err error) {
	var domain string
	base, err := url.Parse(cfApiEndpoint)
	if err != nil {
		log.Error(err)
		return errors.New("error while parsing cloudflare endpoint")
	}
	base.Path += cf.zoneId + "/dns_records?"
	params := url.Values{}
	if cf.hostname == "" {
		domain = cf.domain
	} else {
		domain = cf.hostname + "." + cf.domain
	}
	params.Add("name", domain)
	base.RawQuery = params.Encode()
	endpoint := base.String()
	var data = new(dnsRecord)
	resp := requestMaker(client, header, endpoint)
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return
	}
	if len(data.Errors) > 0 || (len(data.Result) <= 0 && len(data.Errors) >= 0) || data.Result == nil || !data.Success {
		for _, err := range data.Errors {
			log.Error(fmt.Sprintf("%d %s while updating %s ", err.Code, err.Message, domain))
		}
		return errors.New(fmt.Sprintf("Record ID and Record Type not found for %s", domain))
	}
	cf.recordType = data.Result[0].Type
	cf.recordId = data.Result[0].Id
	return nil
}

func (cf *cloudflare) update(client *http.Client, header http.Header) (err error) {
	var domain string
	record := newRecord{
		Kind:    cf.recordType,
		Name:    cf.hostname + "." + cf.domain,
		Content: cf.content,
		Proxied: cf.proxied,
		Ttl:     1,
	}
	base, err := url.Parse(cfApiEndpoint)
	if err != nil {
		log.Error(err)
		return
	}
	base.Path += cf.zoneId + "/dns_records/" + cf.recordId
	params := url.Values{}
	if cf.hostname == "" {
		domain = cf.domain
	} else {
		domain = cf.hostname + "." + cf.domain
	}
	params.Add("name", domain)
	base.RawQuery = params.Encode()
	endpoint := base.String()
	data, err := json.Marshal(record)
	if err != nil {
		log.Error(err)
		return
	}
	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewBuffer(data))
	req.Header = header
	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
			return
		}
	}(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}
	var update = new(updateResponse)
	err = json.Unmarshal(body, &update)
	if err != nil {
		log.Error(err)
		return
	}
	if len(update.Errors) > 0 || !update.Success {
		for _, err := range update.Errors {
			log.Error(fmt.Sprintf("Updating record failed for %s with message %s and error %d: %s",
				domain, update.Messages, err.Code, err.Message))
		}
		return errors.New(fmt.Sprintf("No result while updating %s", domain))
	}
	log.Info(fmt.Sprintf("Successfully updated %s %s", cf.hostname, cf.domain))
	return
}

func UpdateRecord(hostname, domain, content string, proxied bool) {
	record := cloudflare{
		apiToken: os.Getenv("CF_API_TOKEN"),
		hostname: hostname,
		domain:   domain,
		content:  content,
		proxied:  proxied,
	}
	if record.apiToken == "" {
		log.Fatal("Cloudflare Api token not found")
	}
	bearer := "Bearer " + record.apiToken
	headers := http.Header{}
	headers.Add("Authorization", bearer)
	headers.Add("Content-Type", "application/json")
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	err := record.setZoneId(client, headers)
	if err != nil {
		log.Error(err)
		return
	}
	err = record.setRecord(client, headers)
	if err != nil {
		log.Error(err)
		return
	}
	err = record.update(client, headers)
	if err != nil {
		log.Error(err)
		return
	}
	return
}
