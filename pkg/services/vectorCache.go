package services

import (
	"lark-vkm/pkg/openai"
	"time"

	"github.com/patrickmn/go-cache"
)

// 向量查询缓存
type VectorService struct {
	cache *cache.Cache
}

type VectorCacheInterface interface {
	IfProcessed(msg string) *openai.EmbeddingResponse
	SetEmbedings(msg string, response *openai.EmbeddingResponse)
	Clear(msg string) bool
}

var vectorService *VectorService

func (u VectorService) IfProcessed(msg string) *openai.EmbeddingResponse {
	data, found := u.cache.Get(msg)
	if !found {
		return nil
	}
	response, ok := data.(*openai.EmbeddingResponse)
	if !ok {
		return nil
	}
	return response
}

func (u VectorService) SetEmbedings(msg string, response *openai.EmbeddingResponse) {
	u.cache.Set(msg, response, time.Hour*12)
}

func (u VectorService) Clear(msg string) bool {
	u.cache.Delete(msg)
	return true
}

func GetVectorCache() VectorCacheInterface {
	if vectorService == nil {
		vectorService = &VectorService{cache: cache.New(12*time.Hour, 12*time.Hour)}
	}
	return vectorService
}
