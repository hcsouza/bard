package cache

import (
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	. "github.com/hcsouza/bard/config"
	. "github.com/hcsouza/bard/logger"
	"github.com/hcsouza/bard/music"
)

type Client struct {
	memcacheClient *memcache.Client
}

func NewCacheClient() *Client {
	mc := memcache.New(Config.Memcache.Host)
	return &Client{mc}
}

func (client *Client) AddTracksByCountryAndGenre(country, genre string, playlist music.Playlist) error {

	bytes, err := json.Marshal(&playlist)
	if err != nil {
		Logger.Error(fmt.Sprintf("Error on serialize tracks: %s", err))
	}
	itemKey := fmt.Sprintf("%s:%s", country, genre)
	serialized := &memcache.Item{Key: itemKey, Value: []byte(bytes)}

	err = client.memcacheClient.Set(serialized)
	if err != nil {
		Logger.Error(fmt.Sprintf("Error on Set serialized to cache: %s", err))
		return err
	}
	msg := fmt.Sprintf("Tracks for genre: %s on country: %s, adde to cache.", genre, country)
	Logger.Info(msg)
	return err
}

func (client *Client) TracksByCountryAndGenre(country, genre string) (music.Playlist, error) {
	var result music.Playlist
	searchKey := fmt.Sprintf("%s:%s", country, genre)

	it, err := client.memcacheClient.Get(searchKey)
	switch err {
	case memcache.ErrServerError:
		Logger.Error(fmt.Sprintf("Error on Get from memcache: %s ", err))
		return result, err
	case memcache.ErrCacheMiss:
		Logger.Error(fmt.Sprintf("Key not found on cache: %s ", err))
		return result, err
	default:
		if err != nil {
			Logger.Error(fmt.Sprintf("Error on %s ", err))
			return result, err
		}
	}

	err = json.Unmarshal([]byte(it.Value), &result)
	if err != nil {
		Logger.Error(fmt.Sprintf("Error on UnMarshall from memcache: %s ", err))
		return result, err
	}
	msg := fmt.Sprintf("Tracks for genre: %s on country: %s, founded on cache.", genre, country)
	Logger.Info(msg)
	return result, err
}
