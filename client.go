package shplugexporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PlugRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Src     string `json:"src"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

type PlugResponse struct {
	ID     int    `json:"id"`
	Src    string `json:"src"`
	Result struct {
		ID      int     `json:"id"`
		Source  string  `json:"source"`
		Output  bool    `json:"output"`
		Apower  float64 `json:"apower"`
		Voltage float64 `json:"voltage"`
		Freq    float64 `json:"freq"`
		Current float64 `json:"current"`
		Aenergy struct {
			Total    float64   `json:"total"`
			ByMinute []float64 `json:"by_minute"`
			MinuteTs int       `json:"minute_ts"`
		} `json:"aenergy"`
		RetAenergy struct {
			Total    float64   `json:"total"`
			ByMinute []float64 `json:"by_minute"`
			MinuteTs int       `json:"minute_ts"`
		} `json:"ret_aenergy"`
		Temperature struct {
			TC float64 `json:"tC"`
			TF float64 `json:"tF"`
		} `json:"temperature"`
	} `json:"result"`
}

type Client struct {
	device Plug
}

func NewClient(device Plug) *Client {
	return &Client{
		device: device,
	}
}

func (x *Client) SwitchGetStatus() (*PlugResponse, error) {
	req := &PlugRequest{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "Switch.GetStatus",
		Params: map[string]any{
			"id": "0",
		},
	}
	res, err := x.get(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var r *PlugResponse
	bodybytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = json.Unmarshal(bodybytes, &r)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return r, nil
}

func (x *Client) httpClient() *http.Client {
	return http.DefaultClient
}

func (x *Client) get(data *PlugRequest) (*http.Response, error) {
	m, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return x.httpClient().Post(x.device.GetBaseUrl(), "application/json", bytes.NewReader(m))
}
