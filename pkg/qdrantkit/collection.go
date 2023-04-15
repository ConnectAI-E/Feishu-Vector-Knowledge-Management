package qdrantkit

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const (
	collectionApi = "/collections"
)

type CreateCollectionRequest struct {
	Vectors Vectors `json:"vectors"`
}
type Vectors struct {
	Size     int    `json:"size"`
	Distance string `json:"distance"`
}

func (q Qdrant) init() {
	err := q.GetCollection(q.CollectionName)
	if err == nil {
		return
	}
	createCollectionRequest := CreateCollectionRequest{Vectors: Vectors{
		Size:     1536,
		Distance: "Cosine",
	}}
	err = q.CreateCollection(q.CollectionName, createCollectionRequest)
	if err != nil {
		panic("init error:" + err.Error())
	}
}

// CreateCollection creates a collection in qdrant
func (q Qdrant) CreateCollection(name string, createCollectionRequest CreateCollectionRequest) (err error) {
	response := &CommonResponse{}
	var reqBytes []byte
	reqBytes, err = json.Marshal(createCollectionRequest)
	if err != nil {
		log.Println(err.Error())
		return
	}

	body, err := q.Send(http.MethodPut, collectionApi+"/"+name, reqBytes)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	if response.Result == nil {
		return errors.New(response.Status.(map[string]interface{})["error"].(string))
	}
	return

}

func (q Qdrant) GetCollection(name string) (err error) {
	response := &CommonResponse{}

	body, err := q.Send(http.MethodGet, collectionApi+"/"+name, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	if response.Result == nil {
		return errors.New(response.Status.(map[string]interface{})["error"].(string))
	}

	return
}
