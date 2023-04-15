package qdrantkit

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

const (
	pointsApi = "/points"
	searchApi = "/search"
)

type PointRequest struct {
	Points []Point `json:"points"`
}

type Point struct {
	ID      string      `json:"id"`
	Payload interface{} `json:"payload"`
	Vector  []float64   `json:"vector"`
}

type PointResponse struct {
	Result interface{} `json:"result"`
	Status interface{} `json:"status"`
	Time   float64     `json:"time"`
}

type PointSearchRequest struct {
	Params      map[string]interface{} `json:"params"`
	Vector      []float64              `json:"vector"`
	Limit       int                    `json:"limit"`
	WithPayload bool                   `json:"with_payload"`
	WithVector  bool                   `json:"with_vector"`
}
type Match struct {
	Value string `json:"value"`
}
type Must struct {
	Key   string `json:"key"`
	Match Match  `json:"match"`
}

type SearchResult struct {
	ID      string      `json:"id"`
	Version int         `json:"version"`
	Score   float64     `json:"score"`
	Payload interface{} `json:"payload"`
	Vector  []float64   `json:"vector,omitempty"`
}

func (q Qdrant) CreatePoints(collectionName string, pointRequest PointRequest) (err error) {
	response := &CommonResponse{}
	var reqBytes []byte
	reqBytes, err = json.Marshal(pointRequest)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := q.Send(http.MethodPut, collectionApi+"/"+collectionName+pointsApi+"?wait=true", reqBytes)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
		return
	}
	if response.Result == nil {
		log.Println(response.Status)
		return errors.New(response.Status.(map[string]interface{})["error"].(string))
	}
	return

}

func (q Qdrant) SearchPoints(pointSearchRequest PointSearchRequest) (res []SearchResult, err error) {
	response := &CommonResponse{}
	var reqBytes []byte
	reqBytes, err = json.Marshal(pointSearchRequest)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := q.Send(http.MethodPost, collectionApi+"/"+q.CollectionName+pointsApi+searchApi, reqBytes)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
		return
	}
	if response.Result == nil {
		return res, errors.New(response.Status.(map[string]interface{})["error"].(string))
	}
	list := response.Result.([]interface{})
	for _, v := range list {
		sr := SearchResult{}
		err = mapstructure.Decode(v, &sr)
		if err != nil {
			log.Println(err)
			return
		}
		res = append(res, sr)
	}
	return

}
