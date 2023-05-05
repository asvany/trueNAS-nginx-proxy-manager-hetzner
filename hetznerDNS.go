package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	// "golang.org/x/net/ipv4"
)

type HetznerClient struct {
	token       string
	zone_name   string
	zone_id     string
	host_name   string
	required_ip string
	httpClient  *http.Client
}

func NewHetznerClient(token string, zone_name string, host_name string, required_ip string) *HetznerClient {
	client := &HetznerClient{
		token:       token,
		zone_name:   zone_name,
		host_name:   host_name,
		required_ip: required_ip,
		httpClient:  &http.Client{},
	}
	return client
}

func (c *HetznerClient) findZoneID() *HetznerClient {
	// Get Zones (GET https://dns.hetzner.com/api/v1/zones)
	// Create request
	req, err := http.NewRequest("GET", "https://dns.hetzner.com/api/v1/zones", nil)
	if err != nil {
		log.Fatalln("Error on creating request object.\n[ERRO] -", err)
	}

	// Headers
	req.Header.Add("Auth-API-Token", c.token)

	// Fetch Request
	resp, err := c.httpClient.Do(req)

	if err != nil {
		log.Fatalln("Error on response.\n[ERRO] -", err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var data map[string]interface{}

	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		log.Fatalln("Error on unmarshal.\n[ERRO] -", err)
	}
	for _, v := range data["zones"].([]interface{}) {
		if v.(map[string]interface{})["name"] == c.zone_name {
			c.zone_id = v.(map[string]interface{})["id"].(string)
			break
		}
	}

	return c
}

func (c *HetznerClient) findRecordDataByName(name string) (string, string, string) {
	// Get Records (GET https://dns.hetzner.com/api/v1/records?zone_id=)
	// Create request
	req, err := http.NewRequest("GET", "https://dns.hetzner.com/api/v1/records?zone_id="+c.zone_id, nil)
	if err != nil {
		log.Fatalln("Error on creating request object.\n[ERRO] -", err)
	}

	// Headers
	req.Header.Add("Auth-API-Token", c.token)

	// Fetch Request
	resp, err := c.httpClient.Do(req)

	if err != nil {
		log.Fatalln("Error on response.\n[ERRO] -", err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var data map[string]interface{}

	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		log.Fatalln("Error on unmarshal.\n[ERRO] -", err)
	}
	for _, v := range data["records"].([]interface{}) {
		if v.(map[string]interface{})["name"] == name {
			return v.(map[string]interface{})["id"].(string), v.(map[string]interface{})["value"].(string), v.(map[string]interface{})["type"].(string)

		}
	}

	return "", "", ""
}

func (c *HetznerClient) updateRecord(id string, value string, record_type string) {
	// Update Record (PUT https://dns.hetzner.com/api/v1/records/)
	// Create request
	req, err := http.NewRequest("PUT", "https://dns.hetzner.com/api/v1/records/"+id, nil)
	if err != nil {
		log.Fatalln("Error on creating request object.\n[ERRO] -", err)
	}

	// Headers
	req.Header.Add("Auth-API-Token", c.token)
	req.Header.Add("Content-Type", "application/json")

	// Body
	reqBody, err := json.Marshal(map[string]string{
		"value": value,
		"type":  record_type,
	})

	if err != nil {
		log.Fatalln("Error on marshal.\n[ERRO] -", err)
	}

	// Attach body to request
	req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	// Fetch Request
	resp, err := c.httpClient.Do(req)

	if err != nil {
		log.Fatalln("Error on response.\n[ERRO] -", err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var data map[string]interface{}

	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		log.Fatalln("Error on unmarshal.\n[ERRO] -", err)
	}
	fmt.Println(data)
}

func (c *HetznerClient) createRecord(name string, value string, record_type string) {
	// Create Record (POST https://dns.hetzner.com/api/v1/records)
	// Create request

	json := []byte(fmt.Sprintf(`{"value": "%v" ,"ttl": 86400,"type": "%v" ,"name": "%v","zone_id": "%v"}`, value, record_type, name, c.zone_id))
	body := bytes.NewBuffer(json)

	// Create request
	req, err := http.NewRequest("POST", "https://dns.hetzner.com/api/v1/records", body)
	if err != nil {
		log.Fatalln("Error on creating request object.\n[ERRO] -", err)
	}

	// Headers
	req.Header.Add("Auth-API-Token", c.token)
	req.Header.Add("Content-Type", "application/json")

	// Fetch Request
	resp, err := c.httpClient.Do(req)

	if err != nil {
		log.Fatalln("Error on response.\n[ERRO] -", err)
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	// Display Results
	fmt.Println("response Status : ", resp.Status)
	fmt.Println("response Headers : ", resp.Header)
	fmt.Println("response Body : ", string(respBody))
}
