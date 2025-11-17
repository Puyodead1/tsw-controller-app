package tswapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type TSWAPIConfig struct {
	BaseURL    string `example:"http://localhost:31270"`
	CommAPIKey string
}

type TSWAPI struct {
	transport *http.Transport
	client    *http.Client
	Config    TSWAPIConfig
}

func (c *TSWAPI) parseApiResponse(r io.ReadCloser) (map[string]any, error) {
	var data map[string]any
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, err
	}

	if error_code, has_error_code := data["errorCode"]; has_error_code {
		return nil, fmt.Errorf("%s: %s", error_code.(string), data["errorMessage"].(string))
	}

	if _, has_result := data["Result"]; !has_result {
		return nil, fmt.Errorf("invalid_response: Invalid response")
	}

	result := data["Result"].(string)
	if result == "Error" {
		return nil, fmt.Errorf("%s: %s", result, data["Message"].(string))
	}

	return data, nil
}

func (c *TSWAPI) executeTswApiRequest(req *http.Request) (map[string]any, error) {
	if c.Config.CommAPIKey == "" {
		return nil, fmt.Errorf("CommAPIKey has not been configured yet")
	}

	req.Header.Add("DTGCommKey", c.Config.CommAPIKey)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return c.parseApiResponse(resp.Body)
}

func (c *TSWAPI) SetInputValue(control string, value float64) error {
	set_path := fmt.Sprintf("/set/CurrentDrivableActor/%s.InputValue?Value=%f", control, value)
	req_url := fmt.Sprintf("%s%s", c.Config.BaseURL, set_path)
	set_req, _ := http.NewRequest("PATCH", req_url, nil)
	if _, err := c.executeTswApiRequest(set_req); err != nil {
		return err
	}
	return nil
}

func (c *TSWAPI) LoadAPIKey(path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}

	fh, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fh.Close()

	key_bytes, err := io.ReadAll(fh)
	if err != nil {
		return err
	}

	c.Config.CommAPIKey = string(key_bytes)
	return nil
}

func NewTSWAPI(config TSWAPIConfig) *TSWAPI {
	transport := &http.Transport{
		IdleConnTimeout: 120 * time.Second,
	}
	conn := TSWAPI{
		transport: transport,
		client:    &http.Client{Transport: transport, Timeout: 2 * time.Second},
		Config:    config,
	}
	return &conn
}
