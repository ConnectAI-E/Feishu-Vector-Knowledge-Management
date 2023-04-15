package qdrantkit

import (
	"bytes"
	"io"
	"net/http"
)

func (q Qdrant) Send(httpMethod string, suffix string, reqBytes []byte) (body []byte, err error) {
	req, err := http.NewRequest(httpMethod, q.Host+suffix, bytes.NewBuffer(reqBytes))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	return
}

type CommonResponse struct {
	Result interface{} `json:"result"`
	Status interface{} `json:"status"`
	Time   float64     `json:"time"`
}
