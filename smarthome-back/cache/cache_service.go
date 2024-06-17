package cache

import (
	"encoding/json"
	"fmt"
	models "smarthome-back/models/devices"

	"github.com/allegro/bigcache"
)

type CacheService struct {
	Cache *bigcache.BigCache
}

func NewCacheService(cache *bigcache.BigCache) *CacheService {
	return &CacheService{Cache: cache}
}

func (cs *CacheService) GetFromCache(cacheKey string, result interface{}) (bool, error) {
	cachedData, err := cs.Cache.Get(cacheKey)
	if err == nil {
		if err := json.Unmarshal(cachedData, result); err == nil {
			fmt.Println("Data from cache.")
			return true, nil
		}
	}
	return false, err
}

func (cs *CacheService) SetToCache(cacheKey string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return cs.Cache.Set(cacheKey, jsonData)
}

func (cs *CacheService) AddDevicesByRealEstate(estateId int, device models.Device) error {
	cacheKey := fmt.Sprintf("devices_%d", estateId)

	var devices []models.Device
	found, _ := cs.GetFromCache(cacheKey, &devices)
	if found {
		devices := append(devices, device)

		if err := cs.SetToCache(cacheKey, devices); err != nil {
			fmt.Println("Cache error:", err)
			return err
		} else {
			fmt.Println("Saved data in cache.")
		}
	}
	return nil
}
